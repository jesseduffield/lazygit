#!/usr/bin/perl

use strict;
use warnings;
use IPC::Open2;

# An example hook script to integrate Watchman
# (https://facebook.github.io/watchman/) with git to speed up detecting
# new and modified files.
#
# The hook is passed a version (currently 1) and a time in nanoseconds
# formatted as a string and outputs to stdout all files that have been
# modified since the given time. Paths must be relative to the root of
# the working tree and separated by a single NUL.
#
# To enable this hook, rename this file to "query-watchman" and set
# 'git config core.fsmonitor .git/hooks/query-watchman'
#
my ($version, $time) = @ARGV;

# Check the hook interface version

if ($version == 1) {
	# convert nanoseconds to seconds
	$time = int $time / 1000000000;
} else {
	die "Unsupported query-fsmonitor hook version '$version'.\n" .
	    "Falling back to scanning...\n";
}

my $git_work_tree;
if ($^O =~ 'msys' || $^O =~ 'cygwin') {
	$git_work_tree = Win32::GetCwd();
	$git_work_tree =~ tr/\\/\//;
} else {
	require Cwd;
	$git_work_tree = Cwd::cwd();
}

my $retry = 1;

launch_watchman();

sub launch_watchman {

	my $pid = open2(\*CHLD_OUT, \*CHLD_IN, 'watchman -j --no-pretty')
	    or die "open2() failed: $!\n" .
	    "Falling back to scanning...\n";

	# In the query expression below we're asking for names of files that
	# changed since $time but were not transient (ie created after
	# $time but no longer exist).
	#
	# To accomplish this, we're using the "since" generator to use the
	# recency index to select candidate nodes and "fields" to limit the
	# output to file names only. Then we're using the "expression" term to
	# further constrain the results.
	#
	# The category of transient files that we want to ignore will have a
	# creation clock (cclock) newer than $time_t value and will also not
	# currently exist.

	my $query = <<"	END";
		["query", "$git_work_tree", {
			"since": $time,
			"fields": ["name"],
			"expression": ["not", ["allof", ["since", $time, "cclock"], ["not", "exists"]]]
		}]
	END

	print CHLD_IN $query;
	close CHLD_IN;
	my $response = do {local $/; <CHLD_OUT>};

	die "Watchman: command returned no output.\n" .
	    "Falling back to scanning...\n" if $response eq "";
	die "Watchman: command returned invalid output: $response\n" .
	    "Falling back to scanning...\n" unless $response =~ /^\{/;

	my $json_pkg;
	eval {
		require JSON::XS;
		$json_pkg = "JSON::XS";
		1;
	} or do {
		require JSON::PP;
		$json_pkg = "JSON::PP";
	};

	my $o = $json_pkg->new->utf8->decode($response);

	if ($retry > 0 and $o->{error} and $o->{error} =~ m/unable to resolve root .* directory (.*) is not watched/) {
		print STDERR "Adding '$git_work_tree' to watchman's watch list.\n";
		$retry--;
		qx/watchman watch "$git_work_tree"/;
		die "Failed to make watchman watch '$git_work_tree'.\n" .
		    "Falling back to scanning...\n" if $? != 0;

		# Watchman will always return all files on the first query so
		# return the fast "everything is dirty" flag to git and do the
		# Watchman query just to get it over with now so we won't pay
		# the cost in git to look up each individual file.
		print "/\0";
		eval { launch_watchman() };
		exit 0;
	}

	die "Watchman: $o->{error}.\n" .
	    "Falling back to scanning...\n" if $o->{error};

	binmode STDOUT, ":utf8";
	local $, = "\0";
	print @{$o->{files}};
}
