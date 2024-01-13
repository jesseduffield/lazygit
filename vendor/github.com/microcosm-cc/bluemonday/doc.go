// Copyright (c) 2014, David Kitchen <david@buro9.com>
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
//
// * Neither the name of the organisation (Microcosm) nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

/*
Package bluemonday provides a way of describing an allowlist of HTML elements
and attributes as a policy, and for that policy to be applied to untrusted
strings from users that may contain markup. All elements and attributes not on
the allowlist will be stripped.

The default bluemonday.UGCPolicy().Sanitize() turns this:

	Hello <STYLE>.XSS{background-image:url("javascript:alert('XSS')");}</STYLE><A CLASS=XSS></A>World

Into the more harmless:

	Hello World

And it turns this:

	<a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">XSS<a>

Into this:

	XSS

Whilst still allowing this:

	<a href="http://www.google.com/">
	  <img src="https://ssl.gstatic.com/accounts/ui/logo_2x.png"/>
	</a>

To pass through mostly unaltered (it gained a rel="nofollow"):

	<a href="http://www.google.com/" rel="nofollow">
	  <img src="https://ssl.gstatic.com/accounts/ui/logo_2x.png"/>
	</a>

The primary purpose of bluemonday is to take potentially unsafe user generated
content (from things like Markdown, HTML WYSIWYG tools, etc) and make it safe
for you to put on your website.

It protects sites against XSS (http://en.wikipedia.org/wiki/Cross-site_scripting)
and other malicious content that a user interface may deliver. There are many
vectors for an XSS attack (https://www.owasp.org/index.php/XSS_Filter_Evasion_Cheat_Sheet)
and the safest thing to do is to sanitize user input against a known safe list
of HTML elements and attributes.

Note: You should always run bluemonday after any other processing.

If you use blackfriday (https://github.com/russross/blackfriday) or
Pandoc (http://johnmacfarlane.net/pandoc/) then bluemonday should be run after
these steps. This ensures that no insecure HTML is introduced later in your
process.

bluemonday is heavily inspired by both the OWASP Java HTML Sanitizer
(https://code.google.com/p/owasp-java-html-sanitizer/) and the HTML Purifier
(http://htmlpurifier.org/).

We ship two default policies, one is bluemonday.StrictPolicy() and can be
thought of as equivalent to stripping all HTML elements and their attributes as
it has nothing on its allowlist.

The other is bluemonday.UGCPolicy() and allows a broad selection of HTML
elements and attributes that are safe for user generated content. Note that
this policy does not allow iframes, object, embed, styles, script, etc.

The essence of building a policy is to determine which HTML elements and
attributes are considered safe for your scenario. OWASP provide an XSS
prevention cheat sheet ( https://www.google.com/search?q=xss+prevention+cheat+sheet )
to help explain the risks, but essentially:

 1. Avoid allowing anything other than plain HTML elements
 2. Avoid allowing `script`, `style`, `iframe`, `object`, `embed`, `base`
    elements
 3. Avoid allowing anything other than plain HTML elements with simple
    values that you can match to a regexp
*/
package bluemonday
