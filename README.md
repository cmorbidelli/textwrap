Package textwrap is an implementation of Python's textwrap library.
textwrap exports functions for wrapping, indenting, and dedenting
text.  The wrapping functions' behaviors can be customized in a
number of ways, such as dropping leading and trailing whitespace,
expanding tabs, or only wrapping a certain number of lines.  For
use cases demanding a large number of wrapping operations, the
package exports the TextWrapper struct and its methods.

A TextWrapper is initiated with default values, so functions that
use a TextWrapper also accept optional arguments to override some
of the defaults.  For example:

		Wrap("The quick brown fox jumped over the lazy jog",
			Width(10), ExpandTabs(false))

will wrap the text to a line width of 10 without expanding tabs,
but will otherwise use the default behavior for a TextWrapper.

A Textwrapper also contains values defining whitespace and regular
expressions used in the wrapping functions.  These fields are
exported to allow users to customize the functions for non-Latin
character sets.  No options are provided for them because I cannot
guarantee that any changes will produce results that make sense.
