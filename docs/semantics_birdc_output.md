## Semantics of birdc output

### Command `show protocols all`

Output is generated in `/nest/proto.c` (BIRD sourcecode).

#### BGP protocol example (DE-CIX)

In `/nest/proto.c:1476` method `proto_show_stats()` displays information from the statistics struct. All values originate from individual fields in the struct, there is no
redundant storage of information in Bird.  

`birdc show protocols all`
```
R194_129 BGP      T1241_nada_ripe up     2018-06-21 17:42:44  Established
  Description:    Nada & Co.
  Preference:     100
  Input filter:   (unnamed)
  Output filter:  (unnamed)
  Import limit:   200000
    Action:       disable
  Routes:         161 imported, 0 filtered, 164282 exported, 123189 preferred
  Route change stats:     received   rejected   filtered    ignored   accepted
    Import updates:            161          0          0          0        161
    Import withdraws:            0          0        ---          0          0
    Export updates:         226412        322         21        ---     226069
    Export withdraws:           67        ---        ---        ---         67
  BGP state:          Established
    Neighbor address: 172.31.194.129
    Neighbor AS:      1241
    Neighbor ID:      172.31.194.129
    Neighbor caps:    refresh enhanced-refresh AS4
    Session:          external route-server AS4
    Source address:   172.31.192.157
    Route limit:      161/200000
    Hold timer:       118/180
    Keepalive timer:  33/60
```
The meaning of the corresponding fields to values of the birdc output is evident
from the comments after the declarations in `/nest/protocol.h`
```
/* Protocol statistics */
struct proto_stats {
  /* Import - from protocol to core */
  u32 imp_routes;		/* Number of routes successfully imported to the (adjacent) routing table */
  u32 filt_routes;		/* Number of routes rejected in import filter but kept in the routing table */
  u32 pref_routes;		/* Number of routes that are preferred, sum over all routing tables */
  u32 imp_updates_received;	/* Number of route updates received */
  u32 imp_updates_invalid;	/* Number of route updates rejected as invalid */
  u32 imp_updates_filtered;	/* Number of route updates rejected by filters */
  u32 imp_updates_ignored;	/* Number of route updates rejected as already in route table */
  u32 imp_updates_accepted;	/* Number of route updates accepted and imported */
  u32 imp_withdraws_received;	/* Number of route withdraws received */
  u32 imp_withdraws_invalid;	/* Number of route withdraws rejected as invalid */
  u32 imp_withdraws_ignored;	/* Number of route withdraws rejected as already not in route table */
  u32 imp_withdraws_accepted;	/* Number of route withdraws accepted and processed */

  /* Export - from core to protocol */
  u32 exp_routes;		/* Number of routes successfully exported to the protocol */
  u32 exp_updates_received;	/* Number of route updates received */
  u32 exp_updates_rejected;	/* Number of route updates rejected by protocol */
  u32 exp_updates_filtered;	/* Number of route updates rejected by filters */
  u32 exp_updates_accepted;	/* Number of route updates accepted and exported */
  u32 exp_withdraws_received;	/* Number of route withdraws received */
  u32 exp_withdraws_accepted;	/* Number of route withdraws accepted and processed */
};
```

What does the number in `Route limit: %d/%d` mean?
The first number is `stats.imp_routes + stats.filt_routes`, the second is the limit. In the protocol output `stats.imp_routes` is Routes %d imported and `stats.filt_routes` is Routes %d filtered.

The complete protocols of one example neighbor that has multiple routers.
```
M1241_nada_ripe Pipe     master   up     2018-06-21 17:39:31  => T1241_nada_ripe
  Description:    Nada & Co.
  Preference:     70
  Input filter:   in_nada_ripe
  Output filter:  (unnamed)
  Routes:         455 imported, 281241 exported
  Route change stats:     received   rejected   filtered    ignored   accepted
    Import updates:         329913     329425         33          0        455
    Import withdraws:         1249          0        ---          0          0
    Export updates:         486963        455     157083          0     329425
    Export withdraws:         1249          0        ---          0       1249

C1241_nada_ripe Pipe     Collector up     2018-06-21 17:39:31  => T1241_nada_ripe
  Description:    Nada & Co.
  Preference:     70
  Input filter:   in_nada_ripe
  Output filter:  REJECT
  Routes:         455 imported, 0 exported
  Route change stats:     received   rejected   filtered    ignored   accepted
    Import updates:         329913          0     329458          0        455
    Import withdraws:         1249          0        ---          0          0
    Export updates:         440307        455     439852          0          0
    Export withdraws:         1252          0        ---       1252          0

R194_129 BGP      T1241_nada_ripe up     2018-06-21 17:42:44  Established
  Description:    Nada & Co.
  Preference:     100
  Input filter:   (unnamed)
  Output filter:  (unnamed)
  Import limit:   200000
    Action:       disable
  Routes:         161 imported, 0 filtered, 164282 exported, 123189 preferred
  Route change stats:     received   rejected   filtered    ignored   accepted
    Import updates:            161          0          0          0        161
    Import withdraws:            0          0        ---          0          0
    Export updates:         226412        322         21        ---     226069
    Export withdraws:           67        ---        ---        ---         67
  BGP state:          Established
    Neighbor address: 172.31.194.129
    Neighbor AS:      1241
    Neighbor ID:      172.31.194.129
    Neighbor caps:    refresh enhanced-refresh AS4
    Session:          external route-server AS4
    Source address:   172.31.192.157
    Route limit:      161/200000
    Hold timer:       118/180
    Keepalive timer:  33/60

R195_130 BGP      T1241_nada_ripe start  2018-06-21 17:39:31  Passive
  Description:    Nada & Co.
  Preference:     100
  Input filter:   (unnamed)
  Output filter:  (unnamed)
  Import limit:   200000
    Action:       disable
  Routes:         0 imported, 0 filtered, 0 exported, 0 preferred
  Route change stats:     received   rejected   filtered    ignored   accepted
    Import updates:              0          0          0          0          0
    Import withdraws:            0          0        ---          0          0
    Export updates:              0          0          0        ---          0
    Export withdraws:            0        ---        ---        ---          0
  BGP state:          Passive
    Neighbor address: 172.31.195.130
    Neighbor AS:      1241

R193_231 BGP      T1241_nada_ripe up     2018-06-21 17:50:04  Established
  Description:    Nada & Co.
  Preference:     100
  Input filter:   (unnamed)
  Output filter:  (unnamed)
  Import limit:   200000
    Action:       disable
  Routes:         158 imported, 0 filtered, 164091 exported, 8291 preferred
  Route change stats:     received   rejected   filtered    ignored   accepted
    Import updates:            158          0          0          0        158
    Import withdraws:            0          0        ---          0          0
    Export updates:         220336         20        172        ---     220144
    Export withdraws:           67        ---        ---        ---         67
  BGP state:          Established
    Neighbor address: 172.31.193.231
    Neighbor AS:      1241
    Neighbor ID:      172.31.193.231
    Neighbor caps:    refresh enhanced-refresh AS4
    Session:          external route-server AS4
    Source address:   172.31.192.157
    Route limit:      158/200000
    Hold timer:       143/180
    Keepalive timer:  33/60
```

### Command `show route `

Birdc output:  
`198.49.1.0/24      via 172.31.194.42 on eno2 [R194_42 2018-07-27 18:47:27] * (100) [AS715i]`

`next-hop` obvious  
`learnt_from` the IP after `from` (missing in example)  
`gateway` IP address after `via`  
