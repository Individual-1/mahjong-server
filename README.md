* Appengine standard + scale to zero
* No websockets so have to use keepalives + reading responses from user
* How do we identify users? Should we use some sorta signed requests?
* Should we just not bother identifying users at all?
    - firebase auth?
* https://cloud.google.com/appengine/docs/standard/go/how-requests-are-handled
* https://cloud.google.com/appengine/docs/standard/python/how-instances-are-managed#instance_states
* https://cloud.google.com/appengine/docs/standard/go/go-differences#migrating-appengine-sdk
* firebase anonymous authentication userid
* https://cloud.google.com/identity-platform/docs/web/app-engine
* identity platform instead of firebase auth

FCM
* https://medium.com/@rody.davis.jr/how-to-send-push-notifications-on-flutter-web-fcm-b3e64f1e2b76
* On registration have clients generate an FCM id and send it to server
* Server associates firebase auth id, FCM id, and a room
* On each FCM addition to a room, update the group notify for that
* 

Alternative:
* Use Firebase realtime database to have clients pull directly from that
* Anonymous auth to register with a room in the service and then get put onto an allow list for a given room name in db
* https://firebase.google.com/docs/database/security/rules-conditions#structuring_your_database_to_support_authentication_conditions
* https://firebase.google.com/docs/auth/admin/custom-claims#set_and_validate_custom_user_claims_via_the_admin_sdk to set the room name as a claim id
* https://firebase.google.com/docs/database/admin/save-data#go_9
* Use firestore, not realtime database

v3:
* Use Cloud Firestore for DB
* Use FCM to notify groups of users of current game state
* One game per user at a time?
* Store game id in the auth token
* Validate game membership with attempts to change the game state, if invalid then invalidate user's token

server endpoints:

Room related, maybe group?
POST /join {rID=roomID, pw=password} - adds room to user's token and tell them to refresh
POST /leave {rID=roomID} - leave room and remove from all user's db entries and token
POST /create {pw=password} - creates a room and ID, returns id and adds it to creator's token

Tile/game related
POST /discard {tile=tileCode} - discard a given tile on your turn, reject if not your turn or you dont have that tile
POST /take - take the last discarded tile if conditions are met, reject if not
POST /draw {rID=roomID} - returns next tile from the wall, server modifies wallIndex

Schema:

rooms - 
    $rID -
        users -
            $uid1 - wind (as string)
            $uid2 - wind
            $uid3 - wind
            $uid4 - wind
        round - round number
        wall - byte array of tiles
        wallIndex - current index to pull from
        pw - password for the room
        turn - which wind's turn it is
        discard - 
            north - byte array with tiles
            west - same as north
            south - same as north
            east - same as north
        hands -
            north - byte array with tiles
            west - byte array with tiles
            south - byte array with tiles
            east - byte array with tiles
users -
    $uid -
        name - name
        rooms -
            $rID -
                score - current score

cloud function stuff
* https://cloud.google.com/functions/docs/first-go
* https://cloud.google.com/functions/docs/securing/managing-access-iam#allowing_unauthenticated_function_invocation

Misc questions
* How do we handle someone leaving a room, do we kil the room for everyone?
    * Pause game until player requirements met, user can rejoin?
* How do we detect a disconnect?
    * keepalive pings? -> user doesnt check in for X minutes then other users in a room get option to remove player
* Should we worry about hand ordering?
    * Just do it client-side and sanity check the discards on the server side