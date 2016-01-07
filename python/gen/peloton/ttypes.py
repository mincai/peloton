#
# Autogenerated by Thrift Compiler (0.9.3)
#
# DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING
#
#  options string: py:tornado,new_style,dynamic,slots,utf8strings
#

from thrift.Thrift import TType, TMessageType, TException, TApplicationException

from thrift.protocol.TBase import TBase, TExceptionBase, TTransport



class StateAlreadyExists(TExceptionBase):
  """
  Raised when a given state already exists

  Attributes:
   - message
   - existingStates: List of state IDs that already exist
  """

  __slots__ = [ 
    'message',
    'existingStates',
   ]

  thrift_spec = (
    None, # 0
    (1, TType.STRING, 'message', None, None, ), # 1
    (2, TType.LIST, 'existingStates', (TType.STRUCT,(StateIdentity, StateIdentity.thrift_spec)), None, ), # 2
  )

  def __init__(self, message=None, existingStates=None,):
    self.message = message
    self.existingStates = existingStates

  def __str__(self):
    return repr(self)

  def __hash__(self):
    value = 17
    value = (value * 31) ^ hash(self.message)
    value = (value * 31) ^ hash(self.existingStates)
    return value


class InvalidGoalState(TExceptionBase):
  """
  Raised when a given state is not valid

  Attributes:
   - message
   - states: List of invalid goal states
  """

  __slots__ = [ 
    'message',
    'states',
   ]

  thrift_spec = (
    None, # 0
    (1, TType.STRING, 'message', None, None, ), # 1
    (2, TType.LIST, 'states', (TType.STRUCT,(GoalState, GoalState.thrift_spec)), None, ), # 2
  )

  def __init__(self, message=None, states=None,):
    self.message = message
    self.states = states

  def __str__(self):
    return repr(self)

  def __hash__(self):
    value = 17
    value = (value * 31) ^ hash(self.message)
    value = (value * 31) ^ hash(self.states)
    return value


class InternalServerError(TExceptionBase):
  """
  Rasied when a Thrift request invocation throws uncaught exception

  Attributes:
   - message
  """

  __slots__ = [ 
    'message',
   ]

  thrift_spec = (
    None, # 0
    (1, TType.STRING, 'message', None, None, ), # 1
  )

  def __init__(self, message=None,):
    self.message = message

  def __str__(self):
    return repr(self)

  def __hash__(self):
    value = 17
    value = (value * 31) ^ hash(self.message)
    return value


class StateIdentity(TBase):
  """
  Attributes:
   - moduleName: Name of the module which handles the goal state
   - instanceId: Identity of the instance where a goal state will be applied.
  For example, this could be the serivce instance key in UNS format.
   - version: Version number of the goal state which is monolithically increasing.
  Goal state clients and executors can use this to guide against
  race conditions using MVCC.
  """

  __slots__ = [ 
    'moduleName',
    'instanceId',
    'version',
   ]

  thrift_spec = (
    None, # 0
    (1, TType.STRING, 'moduleName', None, None, ), # 1
    (2, TType.STRING, 'instanceId', None, None, ), # 2
    (3, TType.I64, 'version', None, None, ), # 3
  )

  def __init__(self, moduleName=None, instanceId=None, version=None,):
    self.moduleName = moduleName
    self.instanceId = instanceId
    self.version = version

  def __hash__(self):
    value = 17
    value = (value * 31) ^ hash(self.moduleName)
    value = (value * 31) ^ hash(self.instanceId)
    value = (value * 31) ^ hash(self.version)
    return value


class GoalState(TBase):
  """
  Generic goal state definitation for each instance and module.
  The goal state represents the desired state of a module for a service
  instance in a Peloton/Mesos cluster. It is typically set by high-level
  applications like uDeploy upgrade workflow, and will be synced to hosts
  and applied by the specific module via goal state agents.

  Attributes:
   - id: The identity of the goal state
   - updatedAt: The timestamp when the goal state is updated
   - data: The opaque state data that will be handled by the module
   - digest: The digest of the goal state content including fields above
  """

  __slots__ = [ 
    'id',
    'updatedAt',
    'data',
    'digest',
   ]

  thrift_spec = (
    None, # 0
    (1, TType.STRUCT, 'id', (StateIdentity, StateIdentity.thrift_spec), None, ), # 1
    (2, TType.STRING, 'updatedAt', None, None, ), # 2
    (3, TType.STRING, 'data', None, None, ), # 3
    (4, TType.STRING, 'digest', None, None, ), # 4
  )

  def __init__(self, id=None, updatedAt=None, data=None, digest=None,):
    self.id = id
    self.updatedAt = updatedAt
    self.data = data
    self.digest = digest

  def __hash__(self):
    value = 17
    value = (value * 31) ^ hash(self.id)
    value = (value * 31) ^ hash(self.updatedAt)
    value = (value * 31) ^ hash(self.data)
    value = (value * 31) ^ hash(self.digest)
    return value


class ActualState(TBase):
  """
  Generic actual state definitation for each instance and module.
  The actual state represents the actual state of a module for a service
  instance running in a Peloton/Mesos cluster. It is reported by goal
  state agents to the master via periodical updates. If there is no recent
  actual state change, only digest but not data will be present in the request.
  In this case, the actual state update can be treated as a keep alive update.

  Attributes:
   - id: The identity of the actual state
   - updatedAt: The timestamp when the actual state is updated
   - data: The opaque state data that is reported by the module
   - digest: The digest of the goal state content including fields above
  """

  __slots__ = [ 
    'id',
    'updatedAt',
    'data',
    'digest',
   ]

  thrift_spec = (
    None, # 0
    (1, TType.STRUCT, 'id', (StateIdentity, StateIdentity.thrift_spec), None, ), # 1
    (2, TType.STRING, 'updatedAt', None, None, ), # 2
    (3, TType.STRING, 'data', None, None, ), # 3
    (4, TType.STRING, 'digest', None, None, ), # 4
  )

  def __init__(self, id=None, updatedAt=None, data=None, digest=None,):
    self.id = id
    self.updatedAt = updatedAt
    self.data = data
    self.digest = digest

  def __hash__(self):
    value = 17
    value = (value * 31) ^ hash(self.id)
    value = (value * 31) ^ hash(self.updatedAt)
    value = (value * 31) ^ hash(self.data)
    value = (value * 31) ^ hash(self.digest)
    return value


class StateQuery(TBase):
  """
  Query for goal or actual states that match the list of attributes such as
  moduleName, instanceId and minimalVersion

  Attributes:
   - moduleName: Name of the module which handles the goal state
   - instanceId: Identity of the instance where a goal state will be applied.
  For example, this could be the serivce instance key in UNS format.
   - minimalVersion: Minimal version number on which the goal or actual states should be
  greater than or equal to.
   - digest: The digest of the goal or actual state. If matches, the data field will
  be omitted in the return states.
  """

  __slots__ = [ 
    'moduleName',
    'instanceId',
    'minimalVersion',
    'digest',
   ]

  thrift_spec = (
    None, # 0
    (1, TType.STRING, 'moduleName', None, None, ), # 1
    (2, TType.STRING, 'instanceId', None, None, ), # 2
    (3, TType.I64, 'minimalVersion', None, None, ), # 3
    (4, TType.STRING, 'digest', None, None, ), # 4
  )

  def __init__(self, moduleName=None, instanceId=None, minimalVersion=None, digest=None,):
    self.moduleName = moduleName
    self.instanceId = instanceId
    self.minimalVersion = minimalVersion
    self.digest = digest

  def __hash__(self):
    value = 17
    value = (value * 31) ^ hash(self.moduleName)
    value = (value * 31) ^ hash(self.instanceId)
    value = (value * 31) ^ hash(self.minimalVersion)
    value = (value * 31) ^ hash(self.digest)
    return value

