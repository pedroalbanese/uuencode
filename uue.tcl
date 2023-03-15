package require base32

set msg [gets stdin]

set mode [lindex $argv 0]
if { $mode eq "-d" } {
	set text [binary decode uuencode $msg]
	puts $text
} elseif { $mode eq "-e" } {
	set text [binary encode uuencode $msg]
	puts $text
}