#!/usr/bin/env python
from deepdiff import DeepDiff
import datetime
import time

from peloton_client.client import PelotonClient
from peloton_client.pbgen.peloton.api import peloton_pb2 as peloton
from peloton_client.pbgen.peloton.api.job import job_pb2 as job
from peloton_client.pbgen.peloton.api.query import query_pb2 as query
from peloton_client.pbgen.peloton.api.task import task_pb2 as task
from peloton_client.pbgen.peloton.api.respool import respool_pb2 as respool

from color_print import (
    print_okblue,
    print_fail,
)

default_timeout = 60


class PelotonClientHelper(object):
    """
    PelotonClientHelper is using PelotonClient for Peloton operation
    """

    def __init__(self, zk_servers, respool_path=None):
        """
        :param zk_servers: dns address of the physical zk dns
        :type client: PelotonClient
        """
        self.zk_servers = zk_servers
        # Generate PelotonClient
        self.client = PelotonClient(
            name='peloton-client',
            zk_servers=zk_servers,
        )
        if not respool_path:
            return

        # Get resource pool id
        self.respool_id = self.lookup_pool(respool_path)

    def lookup_pool(self, respool_path):
        request = respool.LookupRequest(
            path=respool.ResourcePoolPath(value=respool_path),
        )
        try:
            resp = self.client.respool_svc.LookupResourcePoolID(
                request,
                metadata=self.client.resmgr_metadata,
                timeout=default_timeout,
            )
            return resp.id.value
        except Exception, e:
            print_fail(
                'Failed to get resource pool by path  %s: %s' % (
                    respool_path, str(e)
                )
            )
            raise

    def create_job(self, label, name,
                   num_instance, default_task_config,
                   instance_config=None, **extra):
        """
        :param label: the label value of the job
        :param name: the name of the job
        :param respool_id: the id of the resource pool
        :param num_instance: the number of instance of the job
        :param default_task_config: the default task config of the job
        :param instance_config: instance specific task config
        :param extra: extra information of the job

        :type label: str
        :type name: str
        :type respool_id: str
        :type num_instance: int
        :type default_task_config: task.TaskConfig
        :type instance_config: dict<int, task.TaskConfig>
        :type extra: dict

        :rtypr: job.CreateResponse
        """
        request = job.CreateRequest(
            config=job.JobConfig(
                name=name,
                labels=[
                    peloton.Label(
                        key='cluster_name',
                        value=label,
                    ),
                    peloton.Label(
                        key='module_name',
                        value=name,
                    ),
                ],
                owningTeam=extra.get('owningTeam', 'compute'),
                description=extra.get('description', 'compute task'),
                instanceCount=num_instance,
                defaultConfig=default_task_config,
                instanceConfig=instance_config,
                # sla is required by resmgr
                sla=job.SlaConfig(
                    priority=1,
                    preemptible=True,
                ),
                respoolID=peloton.ResourcePoolID(value=self.respool_id),
            ),
        )

        try:
            resp = self.client.job_svc.Create(
                request,
                metadata=self.client.jobmgr_metadata,
                timeout=default_timeout,
            )
            print_okblue('Create job response : %s' % resp)
            return resp
        except Exception, e:
            print_fail('Exception calling Create job :%s' % str(e))
            raise

    def get_job(self, job_id):
        """
        :param job_id: the id of the job
        :type job_id: str

        :rtype: Response
        """
        request = job.GetRequest(
            id=peloton.JobID(value=job_id),
        )
        try:
            resp = self.client.job_svc.Get(
                request,
                metadata=self.client.jobmgr_metadata,
                timeout=default_timeout,
            )
            return resp
        except Exception, e:
            print_fail('Exception calling Get job :%s' % str(e))
            raise

    def get_jobs_by_label(self, label, name, job_states):
        """
        :param label: the label value of the job
        :param name: the name of the job
        :param job_states: the job status

        :type label: str
        :type name: str
        :type job_states: dict

        :rtype: Response
        """
        request = job.QueryRequest(
            respoolID=peloton.ResourcePoolID(value=self.respool_id),
            spec=job.QuerySpec(
                pagination=query.PaginationSpec(
                    offset=0,
                    limit=100,
                ),
                labels=[
                    peloton.Label(
                        key='cluster_name',
                        value=label,
                    ),
                    peloton.Label(
                        key='module_name',
                        value=name,
                    ),
                ],
                jobStates=job_states,

            ),
        )
        try:
            records = self.client.job_svc.Query(
                request,
                metadata=self.client.jobmgr_metadata,
                timeout=default_timeout,
            ).records
            ids = [record.id.value for record in records]
            return ids

        except Exception, e:
            print_fail('Exception calling Get job :%s' % str(e))
            raise

    def get_job_status(self, job_id):
        resp = self.get_job(job_id)
        return resp.jobInfo.runtime

    def stop_job(self, job_id):
        """
        param job_id: id of the job
        type job_id: str

        rtype: job.StopResponse
        """
        request = task.StopRequest(
            jobId=peloton.JobID(value=job_id),
        )
        try:
            print_okblue("Killing all tasks of Job %s" % job_id)
            resp = self.client.task_svc.Stop(
                request,
                metadata=self.client.jobmgr_metadata,
                timeout=default_timeout,
            )
            return resp
        except Exception, e:
            print_fail('Exception calling List Tasks :%s' % str(e))
            raise

    def get_tasks(self, job_id):
        """
        param job_id: id of the job
        type job_id: str

        rtype: job.ListResponse
        """
        request = task.ListRequest(
            jobId=peloton.JobID(value=job_id),
        )
        try:
            resp = self.client.task_svc.List(
                request,
                metadata=self.client.jobmgr_metadata,
                timeout=default_timeout,
            ).result.value
            return resp
        except Exception, e:
            print_fail('Exception calling List Tasks :%s' % str(e))
            raise

    def monitering(self, job_id, target_status, stable_timeout=120):
        """
        monitering will stop if the job status is not changed in stable_timeout
        or the job status meets the target_status. monitering returns a bool
        value whether the job completedd and meet the target status

        rtype: bool

        """
        if not job_id:
            return
        data = []

        def check_finish(task_stats):
            for k, v in target_status.iteritems():
                if task_stats.get(k, 0) < v[0] or task_stats.get(k, 0) > v[1]:
                    return False
            return True

        stable_timestamp = datetime.datetime.now()
        while datetime.datetime.now() - stable_timestamp < datetime.timedelta(
                seconds=stable_timeout):
            job_runtime = self.get_job_status(job_id)
            task_stats = dict(job_runtime.taskStats)
            data.append(task_stats)
            if check_finish(task_stats):
                break
            if len(data) < 2 or DeepDiff(data[-1], data[-2]):
                # new record is different from previous
                stable_timestamp = datetime.datetime.now()
            time.sleep(5)

        if not check_finish(task_stats):
            return False

        return True