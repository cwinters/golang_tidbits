#!/usr/bin/perl
#
use strict;
use DBI;

my ($file, $database, $user, $password, $host, $port) = @ARGV;

$port ||= 3306;
$host ||= 'localhost';

unless ($user && $password && $host && $port && $database) {
  print usage();
  exit(1);
}
my $dsn = "DBI:mysql:database=$database;host=$host;port=$port";
my $dbh = DBI->connect($dsn, $user, $password);
$dbh->{'mysql_use_result'} = 1;

print STDERR "OK: DB connect\n" if ($dbh->ping());
my $query = read_file($file);
print STDERR "OK: Read query file\n" if ($query);

my $sth = $dbh->prepare($query);
print STDERR "OK: Query execute\n" if ($sth->execute());

my $row_count = 0;
my $row_chunk = 25000;
my $start = time;
my $lap = $start;
my $elapsed = -1;
my $ref = [];

while ($ref = $sth->fetchrow_arrayref) {
  $row_count++;
  print join("\t", @{$ref}), "\n";
  if ($row_count % $row_chunk == 0) {
    $elapsed = time - $lap;
    printf STDERR "%d: %d sec @ %5.2f rows/sec\n", $row_count, $elapsed, $row_chunk/$elapsed;
    $lap = time;
  }
}
$sth->finish();

$elapsed = time - $start;
printf STDERR "DONE - %d in %d sec @ %5.2f rows/sec\n", $row_count, $elapsed, $row_count / $elapsed; 

sub read_file {
  my ($file) = @_;
  local $/ = undef;
  open(IN, $file) || die "Cannot read query from $file: $!";
  my $query = <IN>;
  close(IN);
  return $query;
}

