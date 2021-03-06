## wat-go-sdk
The Watershed Analysis Tool (WAT) golang software development kit (sdk).

# Watershed Analysis Tool
The Watershed Analysis Tool (WAT) is a framework that remotely runs a set of linked plugins. WAT is a framework that provides structure to the declaration of inputs and outputs for a software, facilitates the linking of those inputs and outputs for a series of plugins, and distributes the computation of each plugin within the WAT simulation.
A WAT compute consists of a list of linked plugins represented by a directed acyclic graph (DAG) which defines the computational sequence of the plugins for a given event. WAT also allows for multiple executions of that DAG by allowing the specification of many events to be computed. WAT allows for events to be run in parallel as well as nodes within the DAG (as far as the DAG will allow) to be computed in parallel. 

# Plugins
A plugin is a central idea to WAT. A plugin simply allows for an externally generated software package to be integrated into the WAT. The WAT plugins take an input ModelPayload that defines the resolved paths to resources defined as inputs necessary for compute. A plugin typically is some process model which WAT runs once all inputs have been generated by nodes above the plugin in the DAG. A plugin could be as simple as the generation of a set of random numbers for other plugins to consume, or as complex as a dynamic time-step coupled physics based watershed hydraulics model. 

# Software Development Kit
The software development kit (SDK) provides the essential data structures and a handful of utility services to provide the necessary consistency needed for a framework like WAT. 
