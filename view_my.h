#ifndef _VIEW_MY_H_
#define _VIEW_MY_H_

enum ID { idERROR, idREADY, idDELETE, idDESTROY, idSELECT };

void run(char const *, char const *, int);

void setnfiles(int);

void addfile(char const *);

void reload(char const *);

#endif
