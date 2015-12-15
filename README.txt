Grep-like tool with following features:
 * modern RE syntax (https://code.google.com/p/re2/wiki/Syntax)
 * reads .gz and .bz2 files natively
 * highlights matches
 * match replacements using templates (no need to use sed or other tool with different RE syntax)
 * some basic grep options are supported too (-v, -n, -A/B/C)

Template syntax (when doing replacements):

In the template, a variable is denoted by a substring of the form $name or
${name}, where name is a non-empty sequence of letters, digits, and
underscores. A purely numeric name like $1 refers to the submatch with the
corresponding index; other names refer to capturing parentheses named with the
(?P<name>...) syntax. A reference to an out of range or unmatched index or a
name that is not present in the regular expression is replaced with an empty
string. In the $name form, name is taken to be as long as possible: $1x is equivalent
to ${1x}, not ${1}x, and, $10 is equivalent to ${10}, not ${1}0.
To insert a literal $ in the output, use $$ in the template.
