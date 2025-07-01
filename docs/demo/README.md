# Attila Demo
The Attila demo is designed to run on your local machine. You will need
[Nomad](https://developer.hashicorp.com/nomad/install) and
[Docker](https://docs.docker.com/engine/install/)installed and available locally, as well as the
Attila binary.

### Start Nomad
You will start two Nomad clusters representing the `euw1` and `euw2` regions. Once running, you can
navigate to the Nomad UI's at `http://localhost:4646` and `http://localhost:5646` respectively.
```console
$ nomad agent -dev -config nomad_euw1.hcl
```

```console
$ nomad agent -dev -config nomad_euw2.hcl
```

### Run and Configure Attila
From the root directory of the Attila repository you should first run `make` to generate the Attila
binary. From within this directory, you can then start an Attila server.
```console
$ ../../bin/attila server run -log-level=debug
```

Once the server is running you can first configure the target Nomad regions.
```console
$ ../../bin/attila region create attila_region_euw1.hcl
Name        = euw1
Group       = europe
TLS Enabled = false
Create Time = 2025-06-23T20:18:48+01:00
Update Time = 2025-06-23T20:18:48+01:00

Address                Default
http://localhost:4646  true
```

```console
$ ../../bin/attila region create attila_region_euw2.hcl
Name        = euw2
Group       = europe
TLS Enabled = false
Create Time = 2025-06-23T20:18:58+01:00
Update Time = 2025-06-23T20:18:58+01:00

Address                Default
http://localhost:5646  true
```

Once the regions have been configured, you can use the `attila region shell run <region_name>` to
run a Docker container. It will provide access to the Nomad CLI and has the environment populated
with variables as directed by the stored Attila configuration.

Now that the Nomad regions have been configured, you will need to configure a job registration rule.
When these rules are triggered a two phase filter process happens. Attila will process all available
regions against the boolean `region_filter` parameter. All regions that pass the filter will then be
executed against the `region_picker` selector expression which will pick zero or more regions to
deploy the job to.
```console
$ ../../bin/attila job register rule create attila_job_reg_rule.hcl
Name            = platform_namespace
Region Contexts = namespace
Region Filter   = any(region_namespace, {.Name == "platform"})
Region Picker   = filter(regions, .Group == "europe" )
Create Time     = 2025-06-23T20:51:31+01:00
Update Time     = 2025-06-23T20:51:31+01:00
```

The next and final item to configure is the job registration method. These are how Attila processes
incoming job registrations and decides which registration rules to trigger.
```console
$ ../../bin/attila job register method create attila_job_reg_method.hcl
Name        = platform_namespace
Selector    = Namespace == "platform"
Create Time = 2025-06-23T20:54:58+01:00
Update Time = 2025-06-23T20:54:58+01:00

Rules
- platform_namespace
```

### Plan and Run a Job Registration
With Attila now configured, you can create and run Nomad job registration plans. In the first plan
you will see zero regions selected. This is because neither Nomad region has the `platform`
namespace.
```console
$ ../../bin/attila job register plan create nomad_job.nomad.hcl
ID            = 01JYHPQA96SQ5XWQK3CVJERW49
Num Regions   = 0
Job ID        = example
Job Namespace = platform
```

You can create the namespace in the `euw1` and trigger the creation of a new plan. This time you
will see the plan includes registration to `euw1`.
```console
$ nomad namespace apply -address=http://localhost:4646 platform
```
```console
$ ../../bin/attila job register plan create nomad_job.nomad.hcl
ID            = 01JYHQB5NH3T4R2TJF0909NJFH
Num Regions   = 1
Job ID        = example
Job Namespace = platform

Region "euw1" Plan for Task Group "cache":
Ignored Allocations                 = 0
Placed Allocations                  = 1
Migrated Allocations                = 0
Stopped Allocations                 = 0
In-place Updated Allocations        = 0
Destroyed Allocations               = 0
Canary Allocations                  = 0
Preempted Allocations               = 0
Allocation Placement Failures       = 1
Nodes Evaluated                     = 1
Nodes Exhausted                     = 1
Nodes Available In Datacenter "dc1" = 1
Quotas Exhauted                     = <none>
```

You can then create the namespace in the `euw2` and trigger the creation of a new plan. This plan
will include registration to both regions as they both now have the `platform` namespace.
```console
$ nomad namespace apply -address=http://localhost:5646 platform
```
```console
$ ../../bin/attila job register plan create nomad_job.nomad.hcl
ID            = 01JYHQFRQXN0H793NW04G9C7ZP
Num Regions   = 2
Job ID        = example
Job Namespace = platform

Region "euw1" Plan for Task Group "cache":
Ignored Allocations                 = 0
Placed Allocations                  = 1
Migrated Allocations                = 0
Stopped Allocations                 = 0
In-place Updated Allocations        = 0
Destroyed Allocations               = 0
Canary Allocations                  = 0
Preempted Allocations               = 0
Allocation Placement Failures       = 1
Nodes Evaluated                     = 1
Nodes Exhausted                     = 1
Nodes Available In Datacenter "dc1" = 1
Quotas Exhauted                     = <none>


Region "euw2" Plan for Task Group "cache":
Ignored Allocations                 = 0
Placed Allocations                  = 1
Migrated Allocations                = 0
Stopped Allocations                 = 0
In-place Updated Allocations        = 0
Destroyed Allocations               = 0
Canary Allocations                  = 0
Preempted Allocations               = 0
Allocation Placement Failures       = 1
Nodes Evaluated                     = 1
Nodes Exhausted                     = 1
Nodes Available In Datacenter "dc1" = 1
Quotas Exhauted                     = <none>
```

You can now run the plan which will trigger a deployment of the job two both Nomad regions.
```console
$ ../../bin/attila job register plan run 01JYHQFRQXN0H793NW04G9C7ZP nomad_job.nomad.hcl
ID            = 01JYHQKHZV751SCS7DFW0YBNPH
Num Regions   = 2
Job ID        = example
Job Namespace = platform
Partial Error = <none>

Region "euw1" Run:
Eval ID  = 712259ee-1e35-c1ba-7892-e10602e7cec0
Warnings = <none>
Error    = <none>


Region "euw2" Run:
Eval ID  = 30f205c2-aac3-925d-73d2-abf27b2cd4e2
Warnings = <none>
Error    = <none>
```

Via the Nomad CLI or UI, using the information output above, you can see the job registrations have
been submitted to Nomad.
