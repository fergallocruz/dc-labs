Architecture Document
=====================

The architecture type we use is base on a Client-Server Arquitecture.

The Client makes Post Petitions to the server. The server then reponds with either an error or a validation message.

The Server has an Scheduler and Controller, which gives the tasks to available workers.
If one worker is not available, the task will be given to another.
