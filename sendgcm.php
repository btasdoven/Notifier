<?php


// API access key from Google API's Console
define( 'API_ACCESS_KEY', 'AIzaSyDhdyFnigm2EfKj4LgccjytRYcvUWl6aLA' );


$registrationIds = array("APA91bGLsd_d-pZ6iwoABajxTJ03m0zfs0InI3W5gcmPVZXm7l6nxYCyYAt7CMgqHthcUmRLWQsosLJIBHEBbPLLc-46nxhfejh5ONoVLvlb_clZvPvh7kG7vU3F_vmG4Whk53oPCLxW" );

// prep the bundle
$msg = array
(
    'title'         => 'This is a title2. title2',
	'id'		=> '300001',
    'vibrate'   => 1,
    'sound'     => 1
);

$fields = array
(
    'registration_ids'  => $registrationIds,
    'data'              => $msg
);

$headers = array
(
    'Authorization: key=' . API_ACCESS_KEY,
    'Content-Type: application/json'
);

$ch = curl_init();
curl_setopt( $ch,CURLOPT_URL, 'https://android.googleapis.com/gcm/send' );
curl_setopt( $ch,CURLOPT_POST, true );
curl_setopt( $ch,CURLOPT_HTTPHEADER, $headers );
curl_setopt( $ch,CURLOPT_RETURNTRANSFER, true );
curl_setopt( $ch,CURLOPT_SSL_VERIFYPEER, false );
curl_setopt( $ch,CURLOPT_POSTFIELDS, json_encode( $fields ) );
$result = curl_exec($ch );
curl_close( $ch );

echo $result;
?>
