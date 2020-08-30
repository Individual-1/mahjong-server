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
            $uid1 - true
            $uid2 - true
            $uid3 - true
            $uid4 - true
        round - round number
        wall - byte array of tiles
        wallIndex - current index to pull from
        pw - password for the room
        discard - 
            1 - north, byte array with tiles
            2 - west, same as north
            3 - south, same as north
            4 - east, same as north
users -
    $uid -
        name - name
        rooms -
            $rID -
                tiles - byte array, ordered tiling?
                wind - 1-4, north, west, south, east
                score - current score

cloud function stuff
* https://cloud.google.com/functions/docs/first-go
* https://cloud.google.com/functions/docs/securing/managing-access-iam#allowing_unauthenticated_function_invocation

Misc questions
* How do we handle someone leaving a room, do we kil the room for everyone?
* How do we detect a disconnect?    
    * keepalive pings?