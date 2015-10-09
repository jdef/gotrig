package gotrig

/*
#include <sys/stat.h>
#include <fcntl.h>
#include <errno.h>

#define	FTRIG1_PREFIX    = "ftrig1"
#define	FTRIG1_PREFIXLEN = len(FTRIG1_PREFIX) - 1

int open2 (char const *s, unsigned int flags)
{
  register int r ;
  do
    r = open(s, (int)flags) ;
  while ((r == -1) && (errno == EINTR)) ;
  return r ;
}

int open_read (char const *fn)
{
return open2(fn, O_RDONLY | O_NONBLOCK) ;
}

int open_write (char const *fn)
{
return open2(fn, O_WRONLY | O_NONBLOCK) ;
}

int fd_write (int fd, char const *buf, unsigned int len)
{
  register int r ;
  do r = write(fd, buf, len) ;
  while ((r == -1) && (errno == EINTR)) ;
  return r ;
}

int fd_close (int fd)
{
  for (;;)
  {
    if (!close(fd) || errno == EINPROGRESS) break ;
    if (errno != EINTR) return -1 ;
  }
  return 0 ;
}
*/
import "C"
