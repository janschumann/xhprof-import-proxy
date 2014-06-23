<?php

/**
 * @file
 *
 * This is a simple scripts that imports data to xhprof.io
 */

if (!isset($_POST['http_host']) || !isset($_POST['request_method']) || !isset($_POST['request_uri']) || !isset($_POST['xhprof_data'])) {
  echo "Usage: http://" . $_SERVER['SERVER_NAME'] . $_SERVER['REQUEST_URI'] . ' : POST params: http_host, request_method, request_uri, xhprof_data';
  exit(1);
}

$xhProfData = json_decode($_POST['xhprof_data'], TRUE);

if (!is_array($xhProfData)) {
  echo "ERROR: Invalid (non json) or missing profiling data.";
  exit(1);
}

try {
  $_SERVER['HTTP_HOST'] = $_POST['http_host'];
  $_SERVER['REQUEST_METHOD'] = $_POST['request_method'];
  $_SERVER['REQUEST_URI'] = $_POST['request_uri'];

  require_once __DIR__ . '/xhprof/includes/bootstrap.inc.php';

  $xhprofData = new \ay\xhprof\Data($config['pdo']);
  $xhprofData->save($xhProfData);
}
catch (\Exception $e) {
  echo "ERROR: " . $e->getMessage() . "\n";
  exit(1);
}

echo "Imported '" . $_SERVER['REQUEST_METHOD'] . " " . $_SERVER['HTTP_HOST'] . $_SERVER['REQUEST_URI'] . "'";
exit(0);
