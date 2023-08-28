#ifndef DYNAREADOUT_HEADER_H
#define DYNAREADOUT_HEADER_H
#include "dynareadout/src/key.h"

void keyFileParseGoCallback(char *fileName, size_t lineNumber,
                            char *keywordName, card_t *card, size_t cardIndex,
                            void *userData);
int get_errno();

#endif