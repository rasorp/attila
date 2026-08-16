# Attila Demo

This directory contains a self contained demo for running Attila against two
local Nomad dev agents.

## Prerequisites
- `nomad`
- `docker`
- Go toolchain capable of building this repo

## Build Attila
```console
$ make build
```

## Start the Nomad Regions
From `docs/demo`, start both local Nomad agents in separate terminals.
```console
$ nomad agent -dev -config=nomad_euw1.hcl
$ nomad agent -dev -config=nomad_euw2.hcl
```

## Run and Configure Attila
In a third terminal, start Attila.
```console
$ ../../bin/attila server run --state-memory-enabled
```

Create the two Attila region definitions.
```console
$ ../../bin/attila region create attila_region_euw1.hcl
$ ../../bin/attila region create attila_region_euw2.hcl
$ ../../bin/attila region list
Name  Group   Addresses              TLS Enabled
----  -----   ---------              -----------
euw1  europe  http://localhost:4646  false
euw2  europe  http://localhost:5646  false
```

You can inspect one of the regions as well.
```console
$ ../../bin/attila region get euw1
Name                   = euw1
Group                  = europe
TLS Enabled            = false
Address                Default
http://localhost:4646  true
```

Once the regions have been configured, you can use `attila region shell run
<region_name>` to run a Docker container with the Nomad CLI and the region
credentials populated from Attila.

Now that the Nomad regions have been configured, you will need to configure a
job registration rule. The rule uses a `region_picker` strategy pipeline.
Most of the selection logic lives in expr-lang, with two small helper stages:

1. `filter_expr` keeps only regions where the `platform` namespace exists.
2. `expr` performs arbitrary slice transforms such as filtering by group.
3. `random` shuffles the current slice.
4. `limit` truncates the current slice to `n` entries.

The demo file keeps things deterministic and only uses `filter_expr` and `expr`,
but `random` and `limit` are available for other policies.

```console
$ ../../bin/attila job register rule create attila_job_reg_rule.hcl
Name                     = platform_namespace
Region Contexts          = namespace
Region Picker Strategies = filter_expr, expr
Create Time              = 2025-06-23T20:51:31+01:00
Update Time              = 2025-06-23T20:51:31+01:00
```

The next and final item to configure is the job registration method. These are
how Attila processes incoming job registrations and decides which registration
rules to trigger.
```console
$ ../../bin/attila job register method create attila_job_reg_method.hcl
Name        = platform_namespace
Selector    = Namespace == "platform"
Create Time = 2025-06-23T20:54:58+01:00
Update Time = 2025-06-23T20:54:58+01:00

Rules
- platform_namespace
```

## Plan and Run a Job Registration
With Attila now configured, you can create and run Nomad job registration
plans. In the first plan you will see zero regions selected. This is because
neither Nomad region has the `platform` namespace.
```console
$ ../../bin/attila job register plan create nomad_job.nomad.hcl
ID            = 01JYHPQA96SQ5XWQK3CVJERW49
Num Regions   = 0
Job ID        = example
Job Namespace = platform
```

Create the namespace in `euw1` and generate a new plan. This time you will see
`euw1` included in the registration plan.
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

You can then run the registration using the generated plan which will perform the Nomad job
registration.
```console
$ ../../bin/attila job register plan run 01M04MMT6V257F8RFBYEHW8GJB nomad_job.nomad.hcl
ID            = 01M04MRB1VWFHHKQY4SJA11SAV
Num Regions   = 1
Job ID        = example
Job Namespace = platform
Partial Error = <none>

Region "euw1" Run:
Eval ID  = 4236a659-2a5f-7ec3-bcb7-beb83e75c495
Warnings = <none>
Error    = <none>
```
