#ifndef __SKARNET_LIB_CGO__
#define __SKARNET_LIB_CGO__

int open2 (char const *s, unsigned int flags)
;
int open_read (char const *fn)
;
int open_write (char const *fn)
;
int fd_write (int fd, char const *buf, unsigned int len)
;
int fd_close (int fd)
;

#endif
