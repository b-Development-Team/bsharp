#include <stdbool.h>
#include <stdlib.h>
#include <stdarg.h>
#include <strings.h>
#include <stdio.h>
#include <math.h>
#include <ctype.h>

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

string* string_itoa(int val) {
  int length = snprintf(NULL, 0, "%d", val);
  char* data = malloc(length);
  snprintf(data, length + 1, "%d", val);
  string* s = string_new(data, length);
  return s;
}

string* string_ftoa(float val) {
  int length = snprintf(NULL, 0, "%f", val);
  char* data = malloc(length);
  snprintf(data, length + 1, "%f", val);
  string* s = string_new(data, length);
  s->is_static = false;
  return s;
}

static inline string* string_btoa(bool val) {
  return val ? string_new("true", 4) : string_new("false", 5);
}

// Parse int code
long bsharp_atoi(string* s) {
  // Make null terminated buffer
  char* buf = calloc(s->len + 1, 0);
  memcpy(buf, s->data, s->len);
  long out = strtol(buf, NULL, 10);
  free(buf);
  return out;
}

double bsharp_atof(string* s) {
  // Make null terminated buffer
  char* buf = calloc(s->len + 1, 0);
  memcpy(buf, s->data, s->len);
  double out = strtod(buf, NULL);
  free(buf);
  return out;
}