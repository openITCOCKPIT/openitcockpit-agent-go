package packagemanager

import (
	"strings"
	"testing"
)

// Sample output of the command
// LANG=c rpm -qa --qf "%{NAME} %{EPOCHNUM}:%{VERSION}-%{RELEASE} %{DESCRIPTION}\nEND\n"
var rpmQueryOutput = `gpg-pubkey 0:b86b3716-61e69f29 -----BEGIN PGP PUBLIC KEY BLOCK-----
Version: rpm-4.16.1.3 (NSS-3)

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
1WKL
=jk2t
-----END PGP PUBLIC KEY BLOCK-----

END
tzdata 0:2025c-1.el9 This package contains data files with rules for various timezones around
the world.
END
python3-setuptools-wheel 0:53.0.0-15.el9 A Python wheel of setuptools to use with venv.
END
pcre2-syntax 0:10.40-6.el9 This is a set of manual pages that document a syntax of the regular
expressions implemented by the PCRE2 library.
END
ncurses-base 0:6.2-12.20210508.el9 This package contains descriptions of common terminals. Other terminal
descriptions are included in the ncurses-term package.
END
libreport-filesystem 0:2.15.2-6.el9.alma Filesystem layout for libreport
END
dnf-data 0:4.14.0-31.el9.alma.1 Common data and configuration files for DNF
END
almalinux-gpg-keys 0:9.7-1.el9 This package provides the RPM signature keys for AlmaLinux.
END
almalinux-release 0:9.7-1.el9 AlmaLinux release files.
END
almalinux-repos 0:9.7-1.el9 This package provides the package repository files for AlmaLinux.
END
setup 0:2.13.7-10.el9 The setup package contains a set of important system configuration and
setup files, such as passwd, group, and profile.
END
filesystem 0:3.16-5.el9 The filesystem package is one of the basic packages that is installed
on a Linux system. Filesystem contains the basic directory layout
for a Linux operating system, including the correct permissions for
the directories.
END
basesystem 0:11-13.el9 Basesystem defines the components of a basic Red Hat Enterprise Linux 9 system
(for example, the package installation order to use during bootstrapping).
Basesystem should be in every installation of a system, and it
should never be removed.
END
ncurses-libs 0:6.2-12.20210508.el9 The curses library routines are a terminal-independent method of
updating character screens with reasonable optimization.  The ncurses
(new curses) library is a freely distributable replacement for the
discontinued 4.4 BSD classic curses library.

This package contains the ncurses libraries.
END
bash 0:5.1.8-9.el9 The GNU Bourne Again shell (Bash) is a shell or command language
interpreter that is compatible with the Bourne shell (sh). Bash
incorporates useful features from the Korn shell (ksh) and the C shell
(csh). Most sh scripts can be run by bash without modification.
END
libgcc 0:11.5.0-11.el9.alma.1 This package contains GCC shared support library which is needed
e.g. for exception handling support.
END
glibc-minimal-langpack 0:2.34-231.el9_7.2 This is a Meta package that is used to install minimal language packs.
This package ensures you can use C, POSIX, or C.UTF-8 locales, but
nothing else. It is designed for assembling a minimal system.
END
glibc-common 0:2.34-231.el9_7.2 The glibc-common package includes common binaries for the GNU libc
libraries, as well as national language (locale) support.
END
glibc 0:2.34-231.el9_7.2 The glibc package contains standard libraries which are used by
multiple programs on the system. In order to save disk space and
memory, as well as to make upgrading easier, common system code is
kept in one place and shared between programs. This particular package
contains the most important sets of shared libraries: the standard C
library and the standard math library. Without these two libraries, a
Linux system will not function.
END
zlib 0:1.2.11-40.el9 Zlib is a general-purpose, patent-free, lossless data compression
library which is used by many different programs.
END
xz-libs 0:5.2.5-8.el9_0 Libraries for decoding files compressed with LZMA or XZ utils.
END
bzip2-libs 0:1.0.8-10.el9_5
Libraries for applications using the bzip2 compression format.
END
libzstd 0:1.5.5-1.el9 Zstandard compression shared library.
END
libcap 0:2.48-10.el9 libcap is a library for getting and setting POSIX.1e (formerly POSIX 6)
draft 15 capabilities.
END
libxcrypt 0:4.4.18-3.el9 libxcrypt is a modern library for one-way hashing of passwords.  It
supports a wide variety of both modern and historical hashing methods:
yescrypt, gost-yescrypt, scrypt, bcrypt, sha512crypt, sha256crypt,
md5crypt, SunMD5, sha1crypt, NT, bsdicrypt, bigcrypt, and descrypt.
It provides the traditional Unix crypt and crypt_r interfaces, as well
as a set of extended interfaces pioneered by Openwall Linux, crypt_rn,
crypt_ra, crypt_gensalt, crypt_gensalt_rn, and crypt_gensalt_ra.

libxcrypt is intended to be used by login(1), passwd(1), and other
similar programs; that is, to hash a small number of passwords during
an interactive authentication dialogue with a human. It is not suitable
for use in bulk password-cracking applications, or in any other situation
where speed is more important than careful handling of sensitive data.
However, it is intended to be fast and lightweight enough for use in
servers that must field thousands of login attempts per minute.

This version of the library does not provide the legacy API functions
that have been provided by glibc's libcrypt.so.1.
END
sqlite-libs 0:3.34.1-9.el9_7 This package contains the shared library for sqlite.
END
libuuid 0:2.37.4-21.el9 This is the universally unique ID library, part of util-linux.

The libuuid library generates and parses 128-bit universally unique
id's (UUID's).  A UUID is an identifier that is unique across both
space and time, with respect to the space of all UUIDs.  A UUID can
be used for multiple purposes, from tagging objects with an extremely
short lifetime, to reliably identifying very persistent objects
across a network.

See also the "uuid" package, which is a separate implementation.
END
libgpg-error 0:1.42-5.el9 This is a library that defines common error values for all GnuPG
components.  Among these are GPG, GPGSM, GPGME, GPG-Agent, libgcrypt,
pinentry, SmartCard Daemon and possibly more in the future.
END
libstdc++ 0:11.5.0-11.el9.alma.1 The libstdc++ package contains a rewritten standard compliant GCC Standard
C++ Library.
END
libxml2 0:2.9.13-14.el9_7 This library allows to manipulate XML files. It includes support
to read, modify and write XML and HTML files. There is DTDs support
this includes parsing and validation even with complex DtDs, either
at parse time or later once the document has been modified. The output
can be a simple SAX stream or and in-memory DOM like representations.
In this case one can use the built-in XPath and XPointer implementation
to select sub nodes or ranges. A flexible Input/Output mechanism is
available, with existing HTTP and FTP modules and combined to an
URI library.
END
libattr 0:2.5.1-3.el9 This package contains the libattr.so dynamic library which contains
the extended attribute system calls and library functions.
END
libacl 0:2.3.1-4.el9 This package contains the libacl.so dynamic library which contains
the POSIX 1003.1e draft standard 17 functions for manipulating access
control lists.
END
libsmartcols 0:2.37.4-21.el9 This is library for ls-like terminal programs, part of util-linux.
END
popt 0:1.18-8.el9 Popt is a C library for parsing command line parameters. Popt was
heavily influenced by the getopt() and getopt_long() functions, but
it improves on them by allowing more powerful argument expansion.
Popt can parse arbitrary argv[] style arrays and automatically set
variables based on command line arguments. Popt allows command line
arguments to be aliased via configuration files and includes utility
functions for parsing arbitrary strings into argv[] arrays using
shell-like rules.
END
readline 0:8.1-4.el9 The Readline library provides a set of functions that allow users to
edit command lines. Both Emacs and vi editing modes are available. The
Readline library includes additional functions for maintaining a list
of previously-entered command lines for recalling or editing those
lines, and for performing csh-like history expansion on previous
commands.
END
crypto-policies 0:20250905-1.git377cc42.el9_7 This package provides pre-built configuration files with
cryptographic policies for various cryptographic back-ends,
such as SSL/TLS libraries.
END
libgcrypt 0:1.10.0-11.el9 Libgcrypt is a general purpose crypto library based on the code used
in GNU Privacy Guard.  This is a development version.
END
elfutils-libelf 0:0.193-1.el9.alma.1 The elfutils-libelf package provides a DSO which allows reading and
writing ELF files on a high level.  Third party programs depend on
this package to read internals of ELF files.  The programs of the
elfutils package use it also to generate new ELF files.
END
expat 0:2.5.0-5.el9_7.1 This is expat, the C library for parsing XML, written by James Clark. Expat
is a stream oriented XML parser. This means that you register handlers with
the parser prior to starting the parse. These handlers are called when the
parser discovers the associated structures in the document being parsed. A
start tag is an example of the kind of structures for which you may
register handlers.
END
json-c 0:0.14-11.el9 JSON-C implements a reference counting object model that allows you
to easily construct JSON objects in C, output them as JSON formatted
strings and parse JSON formatted strings back into the C representation
of JSON objects.  It aims to conform to RFC 7159.
END
keyutils-libs 0:1.6.3-1.el9 This package provides a wrapper library for the key management facility system
calls.
END
libcap-ng 0:0.8.2-7.el9 Libcap-ng is a library that makes using posix capabilities easier
END
audit-libs 0:3.1.5-7.el9 The audit-libs package contains the dynamic libraries needed for
applications to use the audit framework.
END
libcom_err 0:1.46.5-8.el9 This is the common error description library, part of e2fsprogs.

libcom_err is an attempt to present a common error-handling mechanism.
END
libtasn1 0:4.16.0-9.el9 A library that provides Abstract Syntax Notation One (ASN.1, as specified
by the X.680 ITU-T recommendation) parsing and structures management, and
Distinguished Encoding Rules (DER, as per X.690) encoding and decoding functions.
END
p11-kit 0:0.25.3-3.el9_5 p11-kit provides a way to load and enumerate PKCS#11 modules, as well
as a standard configuration setup for installing PKCS#11 modules in
such a way that they're discoverable.
END
lua-libs 0:5.4.4-4.el9 This package contains the shared libraries for lua.
END
lz4-libs 0:1.9.3-5.el9 This package contains the libaries for lz4.
END
libassuan 0:2.5.5-3.el9 This is the IPC library used by GnuPG 2, GPGME and a few other
packages.
END
file-libs 0:5.39-16.el9
Libraries for applications using libmagic.
END
alternatives 0:1.24-2.el9 alternatives creates, removes, maintains and displays information about the
symbolic links comprising the alternatives system. It is possible for several
programs fulfilling the same or similar functions to be installed on a single
system at the same time.
END
p11-kit-trust 0:0.25.3-3.el9_5 The p11-kit-trust package contains a system trust PKCS#11 module which
contains certificate anchors and blocklists.
END
gdbm-libs 1:1.23-1.el9 Libraries for the Gdbm GNU database indexing library
END
gmp 1:6.2.0-13.el9 The gmp package contains GNU MP, a library for arbitrary precision
arithmetic, signed integers operations, rational numbers and floating
point numbers. GNU MP is designed for speed, for both small and very
large operands. GNU MP is fast because it uses fullwords as the basic
arithmetic type, it uses fast algorithms, it carefully optimizes
assembly code for many CPUs' most common inner loops, and it generally
emphasizes speed over simplicity/elegance in its operations.

Install the gmp package if you need a fast arbitrary precision
library.
END
libsepol 0:3.6-3.el9 Security-enhanced Linux is a feature of the Linux® kernel and a number
of utilities with enhanced security functionality designed to add
mandatory access controls to Linux.  The Security-enhanced Linux
kernel contains new architectural components originally developed to
improve the security of the Flask operating system. These
architectural components provide general support for the enforcement
of many kinds of mandatory access control policies, including those
based on the concepts of Type Enforcement®, Role-based Access
Control, and Multi-level Security.

libsepol provides an API for the manipulation of SELinux binary policies.
It is used by checkpolicy (the policy compiler) and similar tools, as well
as by programs like load_policy that need to perform specific transformations
on binary policies such as customizing policy boolean settings.
END
libsigsegv 0:2.13-4.el9 This is a library for handling page faults in user mode. A page fault
occurs when a program tries to access to a region of memory that is
currently not available. Catching and handling a page fault is a useful
technique for implementing:
  - pageable virtual memory
  - memory-mapped access to persistent databases
  - generational garbage collectors
  - stack overflow handlers
  - distributed shared memory
END
libunistring 0:0.9.10-15.el9 This portable C library implements Unicode string types in three flavours:
(UTF-8, UTF-16, UTF-32), together with functions for character processing
(names, classifications, properties) and functions for string processing
(iteration, formatted output, width, word breaks, line breaks, normalization,
case folding and regular expressions).
END
libidn2 0:2.3.0-7.el9 Libidn2 is an implementation of the IDNA2008 specifications in RFC
5890, 5891, 5892, 5893 and TR46 for internationalized domain names
(IDN). It is a standalone library, without any dependency on libidn.
END
pcre 0:8.44-4.el9 PCRE, Perl-compatible regular expression, library has its own native API, but
a set of wrapper functions that are based on the POSIX API are also supplied
in the libpcreposix library. Note that this just provides a POSIX calling
interface to PCRE: the regular expressions themselves still follow Perl syntax
and semantics. This package provides support for strings in 8-bit and UTF-8
encodings. Detailed change log is provided by pcre-doc package.
END
grep 0:3.6-5.el9 The GNU versions of commonly used grep utilities. Grep searches through
textual input for lines which contain a match to a specified pattern and then
prints the matching lines. GNU's grep utilities include grep, egrep and fgrep.

GNU grep is needed by many scripts, so it shall be installed on every system.
END
pcre2 0:10.40-6.el9 PCRE2 is a re-working of the original PCRE (Perl-compatible regular
expression) library to provide an entirely new API.

PCRE2 is written in C, and it has its own API. There are three sets of
functions, one for the 8-bit library, which processes strings of bytes, one
for the 16-bit library, which processes strings of 16-bit values, and one for
the 32-bit library, which processes strings of 32-bit values. There are no C++
wrappers. This package provides support for strings in 8-bit and UTF-8
encodings. Install pcre2-utf16 or pcre2-utf32 packages for the other ones.

The distribution does contain a set of C wrapper functions for the 8-bit
library that are based on the POSIX regular expression API (see the pcre2posix
man page). These can be found in a library called libpcre2posix. Note that
this just provides a POSIX calling interface to PCRE2; the regular expressions
themselves still follow Perl syntax and semantics. The POSIX API is
restricted, and does not give full access to all of PCRE2's facilities.
END
libselinux 0:3.6-3.el9 Security-enhanced Linux is a feature of the Linux® kernel and a number
of utilities with enhanced security functionality designed to add
mandatory access controls to Linux.  The Security-enhanced Linux
kernel contains new architectural components originally developed to
improve the security of the Flask operating system. These
architectural components provide general support for the enforcement
of many kinds of mandatory access control policies, including those
based on the concepts of Type Enforcement®, Role-based Access
Control, and Multi-level Security.

libselinux provides an API for SELinux applications to get and set
process and file security contexts and to obtain security policy
decisions.  Required for any applications that use the SELinux API.
END
coreutils-single 0:8.32-39.el9 These are the GNU core utilities,
packaged as a single multicall binary.
END
libblkid 0:2.37.4-21.el9 This is block device identification library, part of util-linux.
END
libmount 0:2.37.4-21.el9 This is the device mounting library, part of util-linux.
END
sed 0:4.8-9.el9 The sed (Stream EDitor) editor is a stream or batch (non-interactive)
editor.  Sed takes text as input, performs an operation or set of
operations on the text and outputs the modified text.  The operations
that sed performs (substitutions, deletions, insertions, etc.) can be
specified in a script file or from the command line.
END
libfdisk 0:2.37.4-21.el9 This is library for fdisk-like programs, part of util-linux.
END
gzip 0:1.12-1.el9 The gzip package contains the popular GNU gzip data compression
program. Gzipped files have a .gz extension.

Gzip should be installed on your system, because it is a
very commonly used data compression program.
END
cracklib 0:2.9.6-27.el9 CrackLib tests passwords to determine whether they match certain
security-oriented characteristics, with the purpose of stopping users
from choosing passwords that are easy to guess. CrackLib performs
several tests on passwords: it tries to generate words from a username
and gecos entry and checks those words against the password; it checks
for simplistic patterns in passwords; and it checks for the password
in a dictionary.

CrackLib is actually a library containing a particular C function
which is used to check the password, as well as other C
functions. CrackLib is not a replacement for a passwd program; it must
be used in conjunction with an existing passwd program.

Install the cracklib package if you need a program to check users'
passwords to see if they are at least minimally secure. If you install
CrackLib, you will also want to install the cracklib-dicts package.
END
cracklib-dicts 0:2.9.6-27.el9 The cracklib-dicts package includes the CrackLib dictionaries.
CrackLib will need to use the dictionary appropriate to your system,
which is normally put in /usr/share/dict/words. Cracklib-dicts also
contains the utilities necessary for the creation of new dictionaries.

If you are installing CrackLib, you should also install cracklib-dicts.
END
findutils 1:4.8.0-7.el9 The findutils package contains programs which will help you locate
files on your system.  The find utility searches through a hierarchy
of directories looking for files which match a certain set of criteria
(such as a file name pattern).  The xargs utility builds and executes
command lines from standard input arguments (usually lists of file
names generated by the find command).

You should install findutils because it includes tools that are very
useful for finding things on your system.
END
ca-certificates 0:2025.2.80_v9.0.305-91.el9 This package contains the set of CA certificates chosen by the
Mozilla Foundation for use with the Internet PKI.
END
openssl-fips-provider 1:3.5.1-4.el9_7 OpenSSL is a toolkit for supporting cryptography. The openssl-fips-provider
package provides the fips.so provider, a cryptography provider that follows
FIPS requirements and provides FIPS approved algorithms.
END
openssl-libs 1:3.5.1-4.el9_7 OpenSSL is a toolkit for supporting cryptography. The openssl-libs
package contains the libraries that are used by various applications which
support cryptographic algorithms and protocols.
END
systemd-libs 0:252-55.el9_7.7.alma.1 Libraries for systemd and udev.
END
util-linux-core 0:2.37.4-21.el9 This is a very basic set of Linux utilities that is necessary on
minimal installations.
END
kmod-libs 0:28-11.el9 The kmod-libs package provides runtime libraries for any application that
wishes to load or unload Linux kernel modules from the running system.
END
libarchive 0:3.5.3-6.el9_6 Libarchive is a programming library that can create and read several different
streaming archive formats, including most popular tar variants, several cpio
formats, and both BSD and GNU ar variants. It can also write shar archives and
read ISO9660 CDROM images and ZIP archives.
END
libevent 0:2.1.12-8.el9_4 The libevent API provides a mechanism to execute a callback function
when a specific event occurs on a file descriptor or after a timeout
has been reached. libevent is meant to replace the asynchronous event
loop found in event driven network servers. An application just needs
to call event_dispatch() and can then add or remove events dynamically
without having to change the event loop.
END
openssl 1:3.5.1-4.el9_7 The OpenSSL toolkit provides support for secure communications between
machines. OpenSSL includes a certificate management tool and shared
libraries which provide various cryptographic algorithms and
protocols.
END
python3-pip-wheel 0:21.3.1-1.el9 A Python wheel of pip to use with venv.
END
python3 0:3.9.25-2.el9_7 Python 3.9 is an accessible, high-level, dynamically typed, interpreted
programming language, designed with an emphasis on code readability.
It includes an extensive standard library, and has a vast ecosystem of
third-party libraries.

The python3 package provides the "python3" executable: the reference
interpreter for the Python language, version 3.
The majority of its standard library is provided in the python3-libs package,
which should be installed automatically along with python3.
The remaining parts of the Python standard library are broken out into the
python3-tkinter and python3-test packages, which may need to be installed
separately.

Documentation for Python is provided in the python3-docs package.

Packages containing additional libraries for Python are generally named with
the "python3-" prefix.
END
python3-libs 0:3.9.25-2.el9_7 This package contains runtime libraries for use by Python:
- the majority of the Python standard library
- a dynamically linked library for use by applications that embed Python as
  a scripting language, and by the main "python3" executable
END
libsemanage 0:3.6-5.el9_6 Security-enhanced Linux is a feature of the Linux® kernel and a number
of utilities with enhanced security functionality designed to add
mandatory access controls to Linux.  The Security-enhanced Linux
kernel contains new architectural components originally developed to
improve the security of the Flask operating system. These
architectural components provide general support for the enforcement
of many kinds of mandatory access control policies, including those
based on the concepts of Type Enforcement®, Role-based Access
Control, and Multi-level Security.

libsemanage provides an API for the manipulation of SELinux binary policies.
It is used by checkpolicy (the policy compiler) and similar tools, as well
as by programs like load_policy that need to perform specific transformations
on binary policies such as customizing policy boolean settings.
END
shadow-utils 2:4.9-15.el9 The shadow-utils package includes the necessary programs for
converting UNIX password files to the shadow password format, plus
programs for managing user and group accounts. The pwconv command
converts passwords to the shadow password format. The pwunconv command
unconverts shadow passwords and generates a passwd file (a standard
UNIX password file). The pwck command checks the integrity of password
and shadow files. The lastlog command prints out the last login times
for all users. The useradd, userdel, and usermod commands are used for
managing user accounts. The groupadd, groupdel, and groupmod commands
are used for managing group accounts.
END
libutempter 0:1.2.1-6.el9 This library provides interface for terminal emulators such as
screen and xterm to record user sessions to utmp and wtmp files.
END
mpfr 0:4.1.0-7.el9 The MPFR library is a C library for multiple-precision floating-point
computations with "correct rounding". The MPFR is efficient and
also has a well-defined semantics. It copies the good ideas from the
ANSI/IEEE-754 standard for double-precision floating-point arithmetic
(53-bit mantissa). MPFR is based on the GMP multiple-precision library.
END
gawk 0:5.1.0-6.el9 The gawk package contains the GNU version of AWK text processing utility. AWK is
a programming language designed for text processing and typically used as a data
extraction and reporting tool.

The gawk utility can be used to do quick and easy text pattern matching,
extracting or reformatting. It is considered to be a standard Linux tool for
text processing.
END
keyutils 0:1.6.3-1.el9 Utilities to control the kernel key management facility and to provide
a mechanism by which the kernel call back to user space to get a key
instantiated.
END
libcomps 0:0.1.18-1.el9 Libcomps is library for structure-like manipulation with content of
comps XML files. Supports read/write XML file, structure(s) modification.
END
python3-libcomps 0:0.1.18-1.el9 Python3 bindings for libcomps library.
END
acl 0:2.3.1-4.el9 This package contains the getfacl and setfacl utilities needed for
manipulating access control lists.
END
attr 0:2.5.1-3.el9 A set of tools for manipulating extended attributes on filesystem
objects, in particular getfattr(1) and setfattr(1).
An attr(1) command is also provided which is largely compatible
with the SGI IRIX tool of the same name.
END
libksba 0:1.5.1-7.el9 KSBA (pronounced Kasbah) is a library to make X.509 certificates as
well as the CMS easily accessible by other applications.  Both
specifications are building blocks of S/MIME and TLS.
END
libdb 0:5.3.28-57.el9_6 The Berkeley Database (Berkeley DB) is a programmatic toolkit that
provides embedded database support for both traditional and
client/server applications. The Berkeley DB includes B+tree, Extended
Linear Hashing, Fixed and Variable-length record access methods,
transactions, locking, logging, shared memory caching, and database
recovery. The Berkeley DB supports C, C++, and Perl APIs. It is
used by many applications, including Python and Perl, so this should
be installed on all systems.
END
libeconf 0:0.4.1-4.el9 libeconf is a highly flexible and configurable library to parse and manage
key=value configuration files. It reads configuration file snippets from
different directories and builds the final configuration file from it.
END
libpwquality 0:1.4.4-8.el9 This is a library for password quality checks and generation
of random passwords that pass the checks.
This library uses the cracklib and cracklib dictionaries
to perform some of the checks.
END
pam 0:1.5.1-26.el9_6 PAM (Pluggable Authentication Modules) is a system security tool that
allows system administrators to set authentication policy without
having to recompile programs that handle authentication.
END
util-linux 0:2.37.4-21.el9 The util-linux package contains a large variety of low-level system
utilities that are necessary for a Linux system to function. Among
others, Util-linux contains the fdisk configuration tool and the login
program.
END
libgomp 0:11.5.0-11.el9.alma.1 This package contains GCC shared support library which is needed
for OpenMP v4.5 support.
END
libnghttp2 0:1.43.0-6.el9 libnghttp2 is a library implementing the Hypertext Transfer Protocol
version 2 (HTTP/2) protocol in C.
END
libseccomp 0:2.5.2-2.el9 The libseccomp library provides an easy to use interface to the Linux Kernel's
syscall filtering mechanism, seccomp.  The libseccomp API allows an application
to specify which syscalls, and optionally which syscall arguments, the
application is allowed to execute, all of which are enforced by the Linux
Kernel.
END
libtool-ltdl 0:2.4.6-46.el9 The libtool-ltdl package contains the GNU Libtool Dynamic Module Loader, a
library that provides a consistent, portable interface which simplifies the
process of using dynamic modules.

These runtime libraries are needed by programs that link directly to the
system-installed ltdl libraries; they are not needed by software built using
the rest of the GNU Autotools (including GNU Autoconf and GNU Automake).
END
libverto 0:0.3.2-3.el9 libverto provides a way for libraries to expose asynchronous interfaces
without having to choose a particular event loop, offloading this
decision to the end application which consumes the library.

If you are packaging an application, not library, based on libverto,
you should depend either on a specific implementation module or you
can depend on the virtual provides 'libverto-module-base'. This will
ensure that you have at least one module installed that provides io,
timeout and signal functionality. Currently glib is the only module
that does not provide these three because it lacks signal. However,
glib will support signal in the future.
END
libcurl-minimal 0:7.76.1-34.el9 This is a replacement of the 'libcurl' package for minimal installations.  It
comes with a limited set of features compared to the 'libcurl' package.  On the
other hand, the package is smaller and requires fewer run-time dependencies to
be installed.
END
curl-minimal 0:7.76.1-34.el9 This is a replacement of the 'curl' package for minimal installations.  It
comes with a limited set of features compared to the 'curl' package.  On the
other hand, the package is smaller and requires fewer run-time dependencies to
be installed.
END
rpm-libs 0:4.16.1.3-39.el9 This package contains the RPM shared libraries.
END
rpm 0:4.16.1.3-39.el9 The RPM Package Manager (RPM) is a powerful command line driven
package management system capable of installing, uninstalling,
verifying, querying, and updating software packages. Each software
package consists of an archive of files along with information about
the package like its version, a description, etc.
END
libsolv 0:0.7.24-3.el9 A free package dependency solver using a satisfiability algorithm. The
library is based on two major, but independent, blocks:

- Using a dictionary approach to store and retrieve package
  and dependency information.

- Using satisfiability, a well known and researched topic, for
  resolving package dependencies.
END
tpm2-tss 0:3.2.3-1.el9 tpm2-tss is a software stack supporting Trusted Platform Module(TPM) 2.0 system
APIs. It sits between TPM driver and applications, providing TPM2.0 specified
APIs for applications to access TPM module through kernel TPM drivers.
END
ima-evm-utils 0:1.6.2-2.el9 The Trusted Computing Group(TCG) run-time Integrity Measurement Architecture
(IMA) maintains a list of hash values of executables and other sensitive
system files, as they are read or executed. These are stored in the file
systems extended attributes. The Extended Verification Module (EVM) prevents
unauthorized changes to these extended attributes on the file system.
ima-evm-utils is used to prepare the file system for these extended attributes.
END
cyrus-sasl-lib 0:2.1.27-22.el9 The cyrus-sasl-lib package contains shared libraries which are needed by
applications which use the Cyrus SASL library.
END
openldap 0:2.6.8-4.el9 OpenLDAP is an open source suite of LDAP (Lightweight Directory Access
Protocol) applications and development tools. LDAP is a set of
protocols for accessing directory services (usually phone book style
information, but other information is possible) over the Internet,
similar to the way DNS (Domain Name System) information is propagated
over the Internet. The openldap package contains configuration files,
libraries, and documentation for OpenLDAP.
END
libyaml 0:0.2.5-7.el9 YAML is a data serialization format designed for human readability and
interaction with scripting languages.  LibYAML is a YAML parser and
emitter written in C.
END
nettle 0:3.10.1-1.el9 Nettle is a cryptographic library that is designed to fit easily in more
or less any context: In crypto toolkits for object-oriented languages
(C++, Python, Pike, ...), in applications like LSH or GNUPG, or even in
kernel space.
END
gnutls 0:3.8.3-9.el9 GnuTLS is a secure communications library implementing the SSL, TLS and DTLS
protocols and technologies around them. It provides a simple C language
application programming interface (API) to access the secure communications
protocols as well as APIs to parse and write X.509, PKCS #12, OpenPGP and
other required structures.
END
glib2 0:2.68.4-18.el9_7 GLib is the low-level core library that forms the basis for projects
such as GTK+ and GNOME. It provides data structure handling for C,
portability wrappers, and interfaces for such runtime functionality
as an event loop, threads, dynamic loading, and an object system.
END
libmodulemd 0:2.13.0-2.el9 C library for manipulating module metadata files.
See https://github.com/fedora-modularity/libmodulemd/blob/master/README.md for
more details.
END
npth 0:1.6-8.el9 nPth is a non-preemptive threads implementation using an API very similar
to the one known from GNU Pth. It has been designed as a replacement of
GNU Pth for non-ancient operating systems. In contrast to GNU Pth is is
based on the system's standard threads implementation. Thus nPth allows
the use of libraries which are not compatible to GNU Pth.
END
gnupg2 0:2.3.3-4.el9 GnuPG is GNU's tool for secure communication and data storage.  It can
be used to encrypt data and to create digital signatures.  It includes
an advanced key management facility and is compliant with the proposed
OpenPGP Internet standard as described in RFC2440 and the S/MIME
standard as described by several RFCs.

GnuPG 2.0 is a newer version of GnuPG with additional support for
S/MIME.  It has a different design philosophy that splits
functionality up into several modules. The S/MIME and smartcard functionality
is provided by the gnupg2-smime package.
END
gpgme 0:1.15.1-6.el9 GnuPG Made Easy (GPGME) is a library designed to make access to GnuPG
easier for applications.  It provides a high-level crypto API for
encryption, decryption, signing, signature verification and key
management.
END
librepo 0:1.14.5-3.el9 A library providing C and Python (libcURL like) API to downloading repository
metadata.
END
libdnf 0:0.69.0-16.el9.alma.1 A Library providing simplified C and Python API to libsolv.
END
python3-libdnf 0:0.69.0-16.el9.alma.1 Python 3 bindings for the libdnf library.
END
python3-hawkey 0:0.69.0-16.el9.alma.1 Python 3 bindings for the hawkey library.
END
python3-gpg 0:1.15.1-6.el9 gpgme bindings for Python 3.
END
rpm-sign-libs 0:4.16.1.3-39.el9 This package contains the RPM shared libraries for signing packages.
END
systemd-rpm-macros 0:252-55.el9_7.7.alma.1 Just the definitions of rpm macros.

See
https://docs.fedoraproject.org/en-US/packaging-guidelines/Scriptlets/#_systemd
for information how to use those macros.
END
dbus 1:1.12.20-8.el9 D-BUS is a system for sending messages between applications. It is
used both for the system-wide message bus service, and as a
per-user-login-session messaging facility.
END
systemd-pam 0:252-55.el9_7.7.alma.1 Systemd PAM module registers the session with systemd-logind.
END
systemd 0:252-55.el9_7.7.alma.1 systemd is a system and service manager that runs as PID 1 and starts
the rest of the system. It provides aggressive parallelization
capabilities, uses socket and D-Bus activation for starting services,
offers on-demand starting of daemons, keeps track of processes using
Linux control groups, maintains mount and automount points, and
implements an elaborate transactional dependency-based service control
logic. systemd supports SysV and LSB init scripts and works as a
replacement for sysvinit. Other parts of this package are a logging daemon,
utilities to control basic system configuration like the hostname,
date, locale, maintain a list of logged-in users, system accounts,
runtime directories and settings, and daemons to manage simple network
configuration, network time synchronization, log forwarding, and name
resolution.
END
dbus-common 1:1.12.20-8.el9 The dbus-common package provides the configuration and setup files for D-Bus
implementations to provide a System and User Message Bus.
END
dbus-broker 0:28-7.el9 dbus-broker is an implementation of a message bus as defined by the D-Bus
specification. Its aim is to provide high performance and reliability, while
keeping compatibility to the D-Bus reference implementation. It is exclusively
written for Linux systems, and makes use of many modern features provided by
recent Linux kernel releases.
END
elfutils-default-yama-scope 0:0.193-1.el9.alma.1 Yama sysctl setting to enable default attach scope settings
enabling programs to use ptrace attach, access to
/proc/PID/{mem,personality,stack,syscall}, and the syscalls
process_vm_readv and process_vm_writev which are used for
interprocess services, communication and introspection
(like synchronisation, signaling, debugging, tracing and
profiling) of processes.
END
elfutils-libs 0:0.193-1.el9.alma.1 The elfutils-libs package contains libraries which implement DWARF, ELF,
and machine-specific ELF handling and process introspection.  These
libraries are used by the programs in the elfutils package.  The
elfutils-devel package enables building other programs using these
libraries.
END
elfutils-debuginfod-client 0:0.193-1.el9.alma.1 The elfutils-debuginfod-client package contains shared libraries
dynamically loaded from -ldw, which use a debuginfod service
to look up debuginfo and associated data. Also includes a
command-line frontend.
END
binutils-gold 0:2.35.2-67.el9_7.1 This package provides the GOLD linker, which can be used as an alternative to
the default binutils linker (ld.bfd).  The GOLD is generally faster than the
BFD linker, and it supports features such as Identical Code Folding and
Incremental linking.  Unfortunately it is not as well maintained as the BFD
linker, and it may become deprecated in the future.


BuildRequires: bison, m4, gcc-c++

BuildRequires: libstdc++-static


BuildRequires: gcc-c++
Conflicts: gcc-c++ < 4.0.0
END
binutils 0:2.35.2-67.el9_7.1 Binutils is a collection of binary utilities, including ar (for
creating, modifying and extracting from archives), as (a family of GNU
assemblers), gprof (for displaying call graph profile data), ld (the
GNU linker), nm (for listing symbols from object files), objcopy (for
copying and translating object files), objdump (for displaying
information from object files), ranlib (for generating an index for
the contents of an archive), readelf (for displaying detailed
information about binary files), size (for listing the section sizes
of an object or archive file), strings (for listing printable strings
from files), strip (for discarding symbols), and addr2line (for
converting addresses to file and line).
END
rpm-build-libs 0:4.16.1.3-39.el9 This package contains the RPM shared libraries for building packages.
END
python3-rpm 0:4.16.1.3-39.el9 The python3-rpm package contains a module that permits applications
written in the Python programming language to use the interface
supplied by RPM Package Manager libraries.

This package should be installed if you want to develop Python 3
programs that will manipulate RPM packages and databases.
END
python3-dnf 0:4.14.0-31.el9.alma.1 Python 3 interface to DNF.
END
dnf 0:4.14.0-31.el9.alma.1 Utility that allows users to manage packages on their systems.
It supports RPMs, modules and comps groups & environments.
END
yum 0:4.14.0-31.el9.alma.1 Utility that allows users to manage packages on their systems.
It supports RPMs, modules and comps groups & environments.
END
iputils 0:20210202-15.el9_7 The iputils package contains basic utilities for monitoring a network,
including ping. The ping command sends a series of ICMP protocol
ECHO_REQUEST packets to a specified network host to discover whether
the target machine is alive and receiving network traffic.
END
crypto-policies-scripts 0:20250905-1.git377cc42.el9_7 This package provides a tool update-crypto-policies, which applies
the policies provided by the crypto-policies package. These can be
either the pre-built policies from the base package or custom policies
defined in simple policy definition files.

The package also provides a tool fips-mode-setup, which can be used
to enable or disable the system FIPS mode.
END
tar 2:1.34-7.el9 The GNU tar program saves many files together in one archive and can
restore individual files (or all of the files) from that archive. Tar
can also be used to add supplemental files to an archive and to update
or list files in the archive. Tar includes multivolume support,
automatic archive compression/decompression, the ability to perform
remote archives, and the ability to perform incremental and full
backups.

If you want to use tar for remote backups, you also need to install
the rmt package on the remote box.
END
vim-minimal 2:8.2.2637-23.el9_7 VIM (VIsual editor iMproved) is an updated and improved version of the
vi editor.  Vi was the first real screen-based editor for UNIX, and is
still very popular.  VIM improves on vi by adding new features:
multiple windows, multi-level undo, block highlighting and more. The
vim-minimal package includes a minimal version of VIM, providing
the commands vi, view, ex, rvi, and rview. NOTE: The online help is
only available when the vim-common package is installed.
END
xz 0:5.2.5-8.el9_0 XZ Utils are an attempt to make LZMA compression easy to use on free (as in
freedom) operating systems. This is achieved by providing tools and libraries
which are similar to use than the equivalents of the most popular existing
compression algorithms.

LZMA is a general purpose compression algorithm designed by Igor Pavlov as
part of 7-Zip. It provides high compression ratio while keeping the
decompression speed fast.
END
hostname 0:3.23-6.el9 This package provides commands which can be used to display the system's
DNS name, and to display or set its hostname or NIS domain name.
END
less 0:590-6.el9 The less utility is a text file browser that resembles more, but has
more capabilities.  Less allows you to move backwards in the file as
well as forwards.  Since less doesn't have to read the entire input file
before it starts, less starts up more quickly than text editors (for
example, vi).

You should install less because it is a basic utility for viewing text
files, and you'll use it frequently.
END
rootfiles 0:8.1-35.el9 The rootfiles package contains basic required files that are placed
in the root user's account.  These files are basically the same
as those in /etc/skel, which are placed in regular
users' home directories.
END
krb5-libs 0:1.21.1-8.el9_6 Kerberos is a network authentication system. The krb5-libs package
contains the shared libraries needed by Kerberos 5. If you are using
Kerberos, you need to install this package.
END`

// LANG=c dnf --quiet check-update
var dnfCheckupdateOutputForRpmTesting = `Failed to set locale, defaulting to C.UTF-8

tar.x86_64                                                                                                               2:1.34-9.el9_7                                                                                                               baseos`

// dnf updateinfo list --available security
// but this output is faked as the systemd had no security updates at the time of writing
var dnfAvailableSecurityUpdatesForRpmTestingFaked = `Last metadata expiration check: 0:52:58 ago on Wed Jan  7 18:13:57 2026.
ALSA-2025:12345 Moderate/Sec.  tar-2:1.34-9.el9_7.x86_64`

func TestRpmManager_getInstalledPackagesWithCancel(t *testing.T) {

	rpm := RpmManager{}

	if !rpm.IsAvailable() {
		t.Skip("Skipping test; rpm is not available on this system.")
	}

	ctx := t.Context()
	output, err := rpm.getInstalledPackagesWithCancel(ctx)
	if err != nil {
		t.Fatalf("Error getting installed packages: %v", err)
	}

	if !strings.HasPrefix(output, "gpg-pubkey") {
		t.Errorf("Unexpected output, got: %s", output)
	}
}

func TestRpmManager_parseRpmOutput(t *testing.T) {

	rpm := RpmManager{}

	packages, err := rpm.parseRpmOutput(rpmQueryOutput)
	if err != nil {
		t.Fatalf("Error parsing dpkg output: %v", err)
	}

	// Table-driven test cases
	tests := []struct {
		name        string
		pkgName     string
		wantVersion string
		wantDesc    string
	}{
		{
			name:        "bash package",
			pkgName:     "bash",
			wantVersion: "5.1.8-9.el9",
			wantDesc:    "The GNU Bourne Again shell",
		},
		{
			name:        "glibc package",
			pkgName:     "glibc",
			wantVersion: "2.34-231.el9_7.2",
			wantDesc:    "The glibc package contains standard libraries which are used by",
		},
		{
			name:        "libuuid package",
			pkgName:     "libuuid",
			wantVersion: "2.37.4-21.el9",
			wantDesc:    "This is the universally unique ID library, part of util",
		},
		{
			name:        "python3-pip-wheel package",
			pkgName:     "python3-pip-wheel",
			wantVersion: "21.3.1-1.el9",
			wantDesc:    "A Python wheel of pip to use with venv.",
		},
		{
			name:        "crypto-policies-scripts package",
			pkgName:     "crypto-policies-scripts",
			wantVersion: "20250905-1.git377cc42.el9_7",
			wantDesc:    "This package provides a tool update-crypto-policies, which applies",
		},
		{
			name:        "dbus package",
			pkgName:     "dbus",
			wantVersion: "1:1.12.20-8.el9",
			wantDesc:    "D-BUS is a system for sending messages between applications",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := false
			for _, pkg := range packages {
				// Assume pkg has Name, Version, Description fields
				if pkg.Name == tt.pkgName {
					found = true
					if pkg.Version != tt.wantVersion {
						t.Errorf("%s: got version %q, want %q", tt.pkgName, pkg.Version, tt.wantVersion)
					}
					if !strings.HasPrefix(pkg.Description, tt.wantDesc) {
						t.Errorf("%s: got description %q, want %q", tt.pkgName, pkg.Description, tt.wantDesc)
					}
				}
			}
			if !found {
				t.Errorf("Package %q not found in parsed output", tt.pkgName)
			}
		})
	}
}

func TestRpmManager_combineRpmOutputWithDnfCheckupdate(t *testing.T) {

	rpm := RpmManager{}
	dnf := DnfManager{}

	installedPackages, err := rpm.parseRpmOutput(rpmQueryOutput)
	if err != nil {
		t.Fatalf("Error parsing dpkg output: %v", err)
	}

	var securityPackages map[string]bool

	updates, err := dnf.parseDnfCheckUpdateOutput(dnfCheckupdateOutputForRpmTesting, installedPackages, securityPackages)
	if err != nil {
		t.Fatalf("Error combining rpm output with dnf check-update: %v", err)
	}

	if len(updates) != 1 {
		t.Fatalf("Expected 1 update, got %d", len(updates))
	}

	wantPkgName := "tar"
	wantCurrentVersion := "2:1.34-7.el9"
	wantAvailableVersion := "2:1.34-9.el9_7"
	wantIsSecurityUpdate := false
	update := updates[0]
	if update.Name != wantPkgName {
		t.Errorf("Package name: got %q, want %q", update.Name, wantPkgName)
	}
	if update.CurrentVersion != wantCurrentVersion {
		t.Errorf("Current version: got %q, want %q", update.CurrentVersion, wantCurrentVersion)
	}
	if update.AvailableVersion != wantAvailableVersion {
		t.Errorf("Available version: got %q, want %q", update.AvailableVersion, wantAvailableVersion)
	}
	if update.IsSecurityUpdate != wantIsSecurityUpdate {
		t.Errorf("IsSecurityUpdate: got %v, want %v", update.IsSecurityUpdate, wantIsSecurityUpdate)
	}
}

func TestRpmManager_combineRpmOutputWithDnfCheckupdateAndSecurity(t *testing.T) {

	rpm := RpmManager{}
	dnf := DnfManager{}

	installedPackages, err := rpm.parseRpmOutput(rpmQueryOutput)
	if err != nil {
		t.Fatalf("Error parsing dpkg output: %v", err)
	}

	securityPackages, err := dnf.parseDnfSecurityUpdateOutput(dnfAvailableSecurityUpdatesForRpmTestingFaked)
	if err != nil {
		t.Fatalf("Error parsing security updates: %v", err)
	}

	updates, err := dnf.parseDnfCheckUpdateOutput(dnfCheckupdateOutputForRpmTesting, installedPackages, securityPackages)
	if err != nil {
		t.Fatalf("Error combining rpm output with dnf check-update: %v", err)
	}

	if len(updates) != 1 {
		t.Fatalf("Expected 1 update, got %d", len(updates))
	}

	wantPkgName := "tar"
	wantCurrentVersion := "2:1.34-7.el9"
	wantAvailableVersion := "2:1.34-9.el9_7"
	wantIsSecurityUpdate := true
	update := updates[0]
	if update.Name != wantPkgName {
		t.Errorf("Package name: got %q, want %q", update.Name, wantPkgName)
	}
	if update.CurrentVersion != wantCurrentVersion {
		t.Errorf("Current version: got %q, want %q", update.CurrentVersion, wantCurrentVersion)
	}
	if update.AvailableVersion != wantAvailableVersion {
		t.Errorf("Available version: got %q, want %q", update.AvailableVersion, wantAvailableVersion)
	}
	if update.IsSecurityUpdate != wantIsSecurityUpdate {
		t.Errorf("IsSecurityUpdate: got %v, want %v", update.IsSecurityUpdate, wantIsSecurityUpdate)
	}
}

func TestRpmManager_removeEpochFromVersionIfZero(t *testing.T) {
	rpm := RpmManager{}

	tests := []struct {
		input    string
		expected string
	}{
		{"0:2.34-231.el9_7.2", "2.34-231.el9_7.2"},
		{"1:5.1.8-9.el9", "1:5.1.8-9.el9"},
		{"2:1.34-9.el9_7", "2:1.34-9.el9_7"},
		{"5.1.8-9.el9", "5.1.8-9.el9"},
		{"0:2025c-1.el9", "2025c-1.el9"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := rpm.removeEpochFromVersionIfZero(tt.input)
			if got != tt.expected {
				t.Errorf("removeEpochFromVersionIfZero(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}
