#ifndef DYNAREADOUT_HEADER_H
#define DYNAREADOUT_HEADER_H

void keyFileParseGoCallback(char *keywordName, card_t *card, size_t cardIndex,
                            void *userData);
int get_errno();

#endif