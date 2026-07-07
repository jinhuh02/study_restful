<?php
$host = "mysql"; //서비스 이름
$db = "testdb";
$user = "myuser";
$pass = "userpassword";

$dsn = "mysql:host=$host;dbname=$db;charset=utf8mb4";

try{
    $pdo = new PDO($dsn, $user, $pass);
    $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
    echo "<h1>MySQL 연결 성공!</h1>";

    // 테스트로 MySQL 버전 확인
    $version = $pdo->query("SELECT VERSION()")->fetchColumn();
    echo "MySQL 버전 : " . $version;
}catch(PDOException $e){
    echo "<h1>연결 실패</h1>";
    echo $e->getMessage();
}
?>