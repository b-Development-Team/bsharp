#include <stdbool.h>
#include <stdlib.h>
#include <stdarg.h>
#include <strings.h>
#include <stdio.h>

typedef struct string {
  char* data;

  bool is_static;

  int len;
  int refs;
} string;

string* string_from_const(char* val) {
  string* out = malloc(sizeof(string));
  out->refs = 1;
  out->is_static = true;
  out->data = val;
  out->len = strlen(val);
  return out;
}

string* string_new(char* data, int len) {
  string* out = malloc(sizeof(string));
  out->refs = 1;
  out->is_static = false;
  out->data = data;
  out->len = len;
  return out;
}

string* string_concat(int n, ...) {
  // Calc len
  va_list ptr;
  va_start(ptr, n);
  int len;
  for (int i = 0; i < n; i++) {
    string* val = va_arg(ptr, string*);
    len += val->len;
  }
  va_end(ptr);
  
  // Actually concat
  char* out = malloc(len);
  int off = 0;
  va_start(ptr, n);
  for (int i = 0; i < n; i++) {
    string* val = va_arg(ptr, string*);
    memcpy(out + off, val->data, val->len);
    off += val->len;
  }
  va_end(ptr);

  // Return
  return string_new(out, len);
}

inline static void string_print(string* val) {
  printf("%.*s\n", val->len, val->data);
}

string* string_slice(string* val, int start, int end) {
  char* out = malloc(end - start);
  memcpy(out, val->data + start, end - start);
  return string_new(out, end - start);
}

void string_free(string* val) {
  if (val == NULL) {
    return;
  }
  val->refs--;
  if (val->refs == 0) {
    if (!val->is_static) {
      free(val->data);
    }
    free(val);
  }
}

