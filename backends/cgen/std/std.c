#include <stdbool.h>
#include <stdlib.h>
#include <stdarg.h>
#include <strings.h>
#include <stdio.h>
#include <math.h>
#include <ctype.h>
#include <sys/time.h>
#include <stdint.h>
#include <stddef.h>

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
  int len = 0;
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

static inline void string_print(string* val) {
  printf("%.*s\n", val->len, val->data);
}

static inline void string_bounds(string* s, int i, const char* pos) {
  if (i < 0 || i >= s->len) {
    fprintf(stderr, "%s: index out of bounds: %d with length %d\n", pos, i, s->len);
    exit(EXIT_FAILURE);
  }
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

string* string_ind(string* val, int index) {
  char* data = malloc(1);
  data[0] = val->data[index];
  return string_new(data, 1);
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
  return val ? string_from_const("true") : string_from_const("false");
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

int string_cmp(string* l, string* r) {
  int len = l->len;
  if (r->len < len) {
    len = r->len;
  }
  return strncmp(l->data, r->data, len);
}

// Time
long normaltime() {
  return time(NULL);
}

long millitime() {
  struct timeval tv;
  gettimeofday(&tv, NULL);
  return tv.tv_sec * 1000L + tv.tv_usec / 1000;
}

long microtime() {
  struct timeval tv;
  gettimeofday(&tv, NULL);
  return tv.tv_sec * 1000000L + tv.tv_usec;
}

long nanotime() {
  struct timespec tv;
  clock_gettime(CLOCK_REALTIME, &tv);
  return tv.tv_sec * 1000000000L + tv.tv_nsec;
}

// Arrays
typedef struct array {
  int off;
  int len;
  int cap;
  void* data;

  int elemsize;
  int refs;
} array;

array* array_new(int elemsize, int cap) {
  array* a = malloc(sizeof(array));
  a->len = 0;
  a->off = 0;
  a->cap = cap;
  a->data = malloc(elemsize * cap);
  a->elemsize = elemsize;
  a->refs = 1;
  return a;
}

// Get pointer to element at index i.
static inline void* array_get(array* a, int i) {
  return a->data + ((a->off + i) * a->elemsize);
}


typedef void (*array_free_fn)(array*);

void array_free(array* a, array_free_fn free_fn) {
  if (a == NULL) {
    return;
  }

  a->refs--;
  if (a->refs == 0) {
    if (free_fn != NULL) {
      a->len += a->off;
      a->off = 0;
      free_fn(a);
    }
    free(a->data);
    free(a);
  }
}

static inline void array_bounds(array* a, int i, const char* pos) {
  if (i < 0 || i >= a->len) {
    fprintf(stderr, "%s: index out of bounds: %d with length %d\n", pos, i, a->len);
    exit(EXIT_FAILURE);
  }
}

static inline void bsp_panic(string* msg, const char* pos) {
  fprintf(stderr, "%s: %.*s\n", pos, msg->len, msg->data);
  exit(EXIT_FAILURE);
}

void array_grow(array* a, int len) {
  if (a->cap < (a->off + len)) {
    a->cap = len + a->off;
    a->data = realloc(a->data, a->elemsize * a->cap);
  }
}

void array_append(array* a, void* val) {
  array_grow(a, a->len + 1);
  memcpy(array_get(a, a->len), val, a->elemsize);
  a->len++;
}

void array_set(array* a, int i, void* val) {
  memcpy(array_get(a, i), val, a->elemsize);
}

void array_slice(array* a, int start, int end) {
  a->off = start;
  a->len = end - start;
}

// https://github.com/tidwall/hashmap.c

// Copyright 2020 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

static void *(*_malloc)(size_t) = NULL;
static void *(*_realloc)(void *, size_t) = NULL;
static void (*_free)(void *) = NULL;

// hashmap_set_allocator allows for configuring a custom allocator for
// all hashmap library operations. This function, if needed, should be called
// only once at startup and a prior to calling hashmap_new().
void hashmap_set_allocator(void *(*malloc)(size_t), void (*free)(void*)) 
{
    _malloc = malloc;
    _free = free;
}

#define panic(_msg_) { \
    fprintf(stderr, "panic: %s (%s:%d)\n", (_msg_), __FILE__, __LINE__); \
    exit(1); \
}

struct bucket {
    uint64_t hash:48;
    uint64_t dib:16;
};

// hashmap is an open addressed hash map using robinhood hashing.
struct hashmap {
    void *(*malloc)(size_t);
    void *(*realloc)(void *, size_t);
    void (*free)(void *);
    bool oom;
    size_t elsize;
    size_t cap;
    uint64_t seed0;
    uint64_t seed1;
    uint64_t (*hash)(const void *item, uint64_t seed0, uint64_t seed1);
    int (*compare)(const void *a, const void *b, void *udata);
    void (*elfree)(void *item);
    void *udata;
    size_t bucketsz;
    size_t nbuckets;
    size_t count;
    size_t mask;
    size_t growat;
    size_t shrinkat;
    void *buckets;
    void *spare;
    void *edata;
};

static struct bucket *bucket_at(struct hashmap *map, size_t index) {
    return (struct bucket*)(((char*)map->buckets)+(map->bucketsz*index));
}

static void *bucket_item(struct bucket *entry) {
    return ((char*)entry)+sizeof(struct bucket);
}

static uint64_t get_hash(struct hashmap *map, const void *key) {
    return map->hash(key, map->seed0, map->seed1) << 16 >> 16;
}

// hashmap_new_with_allocator returns a new hash map using a custom allocator.
// See hashmap_new for more information information
struct hashmap *hashmap_new_with_allocator(
                            void *(*_malloc)(size_t), 
                            void *(*_realloc)(void*, size_t), 
                            void (*_free)(void*),
                            size_t elsize, size_t cap, 
                            uint64_t seed0, uint64_t seed1,
                            uint64_t (*hash)(const void *item, 
                                             uint64_t seed0, uint64_t seed1),
                            int (*compare)(const void *a, const void *b, 
                                           void *udata),
                            void (*elfree)(void *item),
                            void *udata)
{
    _malloc = _malloc ? _malloc : malloc;
    _realloc = _realloc ? _realloc : realloc;
    _free = _free ? _free : free;
    int ncap = 16;
    if (cap < ncap) {
        cap = ncap;
    } else {
        while (ncap < cap) {
            ncap *= 2;
        }
        cap = ncap;
    }
    size_t bucketsz = sizeof(struct bucket) + elsize;
    while (bucketsz & (sizeof(uintptr_t)-1)) {
        bucketsz++;
    }
    // hashmap + spare + edata
    size_t size = sizeof(struct hashmap)+bucketsz*2;
    struct hashmap *map = _malloc(size);
    if (!map) {
        return NULL;
    }
    memset(map, 0, sizeof(struct hashmap));
    map->elsize = elsize;
    map->bucketsz = bucketsz;
    map->seed0 = seed0;
    map->seed1 = seed1;
    map->hash = hash;
    map->compare = compare;
    map->elfree = elfree;
    map->udata = udata;
    map->spare = ((char*)map)+sizeof(struct hashmap);
    map->edata = (char*)map->spare+bucketsz;
    map->cap = cap;
    map->nbuckets = cap;
    map->mask = map->nbuckets-1;
    map->buckets = _malloc(map->bucketsz*map->nbuckets);
    if (!map->buckets) {
        _free(map);
        return NULL;
    }
    memset(map->buckets, 0, map->bucketsz*map->nbuckets);
    map->growat = map->nbuckets*0.75;
    map->shrinkat = map->nbuckets*0.10;
    map->malloc = _malloc;
    map->realloc = _realloc;
    map->free = _free;
    return map;  
}


// hashmap_new returns a new hash map. 
// Param `elsize` is the size of each element in the tree. Every element that
// is inserted, deleted, or retrieved will be this size.
// Param `cap` is the default lower capacity of the hashmap. Setting this to
// zero will default to 16.
// Params `seed0` and `seed1` are optional seed values that are passed to the 
// following `hash` function. These can be any value you wish but it's often 
// best to use randomly generated values.
// Param `hash` is a function that generates a hash value for an item. It's
// important that you provide a good hash function, otherwise it will perform
// poorly or be vulnerable to Denial-of-service attacks. This implementation
// comes with two helper functions `hashmap_sip()` and `hashmap_murmur()`.
// Param `compare` is a function that compares items in the tree. See the 
// qsort stdlib function for an example of how this function works.
// The hashmap must be freed with hashmap_free(). 
// Param `elfree` is a function that frees a specific item. This should be NULL
// unless you're storing some kind of reference data in the hash.
struct hashmap *hashmap_new(size_t elsize, size_t cap, 
                            uint64_t seed0, uint64_t seed1,
                            uint64_t (*hash)(const void *item, 
                                             uint64_t seed0, uint64_t seed1),
                            int (*compare)(const void *a, const void *b, 
                                           void *udata),
                            void (*elfree)(void *item),
                            void *udata)
{
    return hashmap_new_with_allocator(
        (_malloc?_malloc:malloc),
        (_realloc?_realloc:realloc),
        (_free?_free:free),
        elsize, cap, seed0, seed1, hash, compare, elfree, udata
    );
}

static void free_elements(struct hashmap *map) {
    if (map->elfree) {
        for (size_t i = 0; i < map->nbuckets; i++) {
            struct bucket *bucket = bucket_at(map, i);
            if (bucket->dib) map->elfree(bucket_item(bucket));
        }
    }
}


// hashmap_clear quickly clears the map. 
// Every item is called with the element-freeing function given in hashmap_new,
// if present, to free any data referenced in the elements of the hashmap.
// When the update_cap is provided, the map's capacity will be updated to match
// the currently number of allocated buckets. This is an optimization to ensure
// that this operation does not perform any allocations.
void hashmap_clear(struct hashmap *map, bool update_cap) {
    map->count = 0;
    free_elements(map);
    if (update_cap) {
        map->cap = map->nbuckets;
    } else if (map->nbuckets != map->cap) {
        void *new_buckets = map->malloc(map->bucketsz*map->cap);
        if (new_buckets) {
            map->free(map->buckets);
            map->buckets = new_buckets;
        }
        map->nbuckets = map->cap;
    }
    memset(map->buckets, 0, map->bucketsz*map->nbuckets);
    map->mask = map->nbuckets-1;
    map->growat = map->nbuckets*0.75;
    map->shrinkat = map->nbuckets*0.10;
}


static bool resize(struct hashmap *map, size_t new_cap) {
    struct hashmap *map2 = hashmap_new(map->elsize, new_cap, map->seed1, 
                                       map->seed1, map->hash, map->compare,
                                       map->elfree, map->udata);
    if (!map2) {
        return false;
    }
    for (size_t i = 0; i < map->nbuckets; i++) {
        struct bucket *entry = bucket_at(map, i);
        if (!entry->dib) {
            continue;
        }
        entry->dib = 1;
        size_t j = entry->hash & map2->mask;
        for (;;) {
            struct bucket *bucket = bucket_at(map2, j);
            if (bucket->dib == 0) {
                memcpy(bucket, entry, map->bucketsz);
                break;
            }
            if (bucket->dib < entry->dib) {
                memcpy(map2->spare, bucket, map->bucketsz);
                memcpy(bucket, entry, map->bucketsz);
                memcpy(entry, map2->spare, map->bucketsz);
            }
            j = (j + 1) & map2->mask;
            entry->dib += 1;
        }
	}
    map->free(map->buckets);
    map->buckets = map2->buckets;
    map->nbuckets = map2->nbuckets;
    map->mask = map2->mask;
    map->growat = map2->growat;
    map->shrinkat = map2->shrinkat;
    map->free(map2);
    return true;
}

// hashmap_set inserts or replaces an item in the hash map. If an item is
// replaced then it is returned otherwise NULL is returned. This operation
// may allocate memory. If the system is unable to allocate additional
// memory then NULL is returned and hashmap_oom() returns true.
void *hashmap_set(struct hashmap *map, void *item) {
    if (!item) {
        panic("item is null");
    }
    map->oom = false;
    if (map->count == map->growat) {
        if (!resize(map, map->nbuckets*2)) {
            map->oom = true;
            return NULL;
        }
    }

    
    struct bucket *entry = map->edata;
    entry->hash = get_hash(map, item);
    entry->dib = 1;
    memcpy(bucket_item(entry), item, map->elsize);
    
    size_t i = entry->hash & map->mask;
	for (;;) {
        struct bucket *bucket = bucket_at(map, i);
        if (bucket->dib == 0) {
            memcpy(bucket, entry, map->bucketsz);
            map->count++;
			return NULL;
		}
        if (entry->hash == bucket->hash && 
            map->compare(bucket_item(entry), bucket_item(bucket), 
                         map->udata) == 0)
        {
            memcpy(map->spare, bucket_item(bucket), map->elsize);
            memcpy(bucket_item(bucket), bucket_item(entry), map->elsize);
            return map->spare;
		}
        if (bucket->dib < entry->dib) {
            memcpy(map->spare, bucket, map->bucketsz);
            memcpy(bucket, entry, map->bucketsz);
            memcpy(entry, map->spare, map->bucketsz);
		}
		i = (i + 1) & map->mask;
        entry->dib += 1;
	}
}

// hashmap_get returns the item based on the provided key. If the item is not
// found then NULL is returned.
void *hashmap_get(struct hashmap *map, const void *key) {
    if (!key) {
        panic("key is null");
    }
    uint64_t hash = get_hash(map, key);
	size_t i = hash & map->mask;
	for (;;) {
        struct bucket *bucket = bucket_at(map, i);
		if (!bucket->dib) {
			return NULL;
		}
		if (bucket->hash == hash && 
            map->compare(key, bucket_item(bucket), map->udata) == 0)
        {
            return bucket_item(bucket);
		}
		i = (i + 1) & map->mask;
	}
}

// hashmap_probe returns the item in the bucket at position or NULL if an item
// is not set for that bucket. The position is 'moduloed' by the number of 
// buckets in the hashmap.
void *hashmap_probe(struct hashmap *map, uint64_t position) {
    size_t i = position & map->mask;
    struct bucket *bucket = bucket_at(map, i);
    if (!bucket->dib) {
		return NULL;
	}
    return bucket_item(bucket);
}


// hashmap_delete removes an item from the hash map and returns it. If the
// item is not found then NULL is returned.
void *hashmap_delete(struct hashmap *map, void *key) {
    if (!key) {
        panic("key is null");
    }
    map->oom = false;
    uint64_t hash = get_hash(map, key);
	size_t i = hash & map->mask;
	for (;;) {
        struct bucket *bucket = bucket_at(map, i);
		if (!bucket->dib) {
			return NULL;
		}
		if (bucket->hash == hash && 
            map->compare(key, bucket_item(bucket), map->udata) == 0)
        {
            memcpy(map->spare, bucket_item(bucket), map->elsize);
            bucket->dib = 0;
            for (;;) {
                struct bucket *prev = bucket;
                i = (i + 1) & map->mask;
                bucket = bucket_at(map, i);
                if (bucket->dib <= 1) {
                    prev->dib = 0;
                    break;
                }
                memcpy(prev, bucket, map->bucketsz);
                prev->dib--;
            }
            map->count--;
            if (map->nbuckets > map->cap && map->count <= map->shrinkat) {
                // Ignore the return value. It's ok for the resize operation to
                // fail to allocate enough memory because a shrink operation
                // does not change the integrity of the data.
                resize(map, map->nbuckets/2);
            }
			return map->spare;
		}
		i = (i + 1) & map->mask;
	}
}

// hashmap_count returns the number of items in the hash map.
size_t hashmap_count(struct hashmap *map) {
    return map->count;
}

// hashmap_free frees the hash map
// Every item is called with the element-freeing function given in hashmap_new,
// if present, to free any data referenced in the elements of the hashmap.
void hashmap_free(struct hashmap *map) {
    if (!map) return;
    free_elements(map);
    map->free(map->buckets);
    map->free(map);
}

// hashmap_oom returns true if the last hashmap_set() call failed due to the 
// system being out of memory.
bool hashmap_oom(struct hashmap *map) {
    return map->oom;
}

// hashmap_scan iterates over all items in the hash map
// Param `iter` can return false to stop iteration early.
// Returns false if the iteration has been stopped early.
bool hashmap_scan(struct hashmap *map, 
                  bool (*iter)(const void *item, void *udata), void *udata)
{
    for (size_t i = 0; i < map->nbuckets; i++) {
        struct bucket *bucket = bucket_at(map, i);
        if (bucket->dib) {
            if (!iter(bucket_item(bucket), udata)) {
                return false;
            }
        }
    }
    return true;
}

//-----------------------------------------------------------------------------
// SipHash reference C implementation
//
// Copyright (c) 2012-2016 Jean-Philippe Aumasson
// <jeanphilippe.aumasson@gmail.com>
// Copyright (c) 2012-2014 Daniel J. Bernstein <djb@cr.yp.to>
//
// To the extent possible under law, the author(s) have dedicated all copyright
// and related and neighboring rights to this software to the public domain
// worldwide. This software is distributed without any warranty.
//
// You should have received a copy of the CC0 Public Domain Dedication along
// with this software. If not, see
// <http://creativecommons.org/publicdomain/zero/1.0/>.
//
// default: SipHash-2-4
//-----------------------------------------------------------------------------
static uint64_t SIP64(const uint8_t *in, const size_t inlen, 
                      uint64_t seed0, uint64_t seed1) 
{
#define U8TO64_LE(p) \
    {  (((uint64_t)((p)[0])) | ((uint64_t)((p)[1]) << 8) | \
        ((uint64_t)((p)[2]) << 16) | ((uint64_t)((p)[3]) << 24) | \
        ((uint64_t)((p)[4]) << 32) | ((uint64_t)((p)[5]) << 40) | \
        ((uint64_t)((p)[6]) << 48) | ((uint64_t)((p)[7]) << 56)) }
#define U64TO8_LE(p, v) \
    { U32TO8_LE((p), (uint32_t)((v))); \
      U32TO8_LE((p) + 4, (uint32_t)((v) >> 32)); }
#define U32TO8_LE(p, v) \
    { (p)[0] = (uint8_t)((v)); \
      (p)[1] = (uint8_t)((v) >> 8); \
      (p)[2] = (uint8_t)((v) >> 16); \
      (p)[3] = (uint8_t)((v) >> 24); }
#define ROTL(x, b) (uint64_t)(((x) << (b)) | ((x) >> (64 - (b))))
#define SIPROUND \
    { v0 += v1; v1 = ROTL(v1, 13); \
      v1 ^= v0; v0 = ROTL(v0, 32); \
      v2 += v3; v3 = ROTL(v3, 16); \
      v3 ^= v2; \
      v0 += v3; v3 = ROTL(v3, 21); \
      v3 ^= v0; \
      v2 += v1; v1 = ROTL(v1, 17); \
      v1 ^= v2; v2 = ROTL(v2, 32); }
    uint64_t k0 = U8TO64_LE((uint8_t*)&seed0);
    uint64_t k1 = U8TO64_LE((uint8_t*)&seed1);
    uint64_t v3 = UINT64_C(0x7465646279746573) ^ k1;
    uint64_t v2 = UINT64_C(0x6c7967656e657261) ^ k0;
    uint64_t v1 = UINT64_C(0x646f72616e646f6d) ^ k1;
    uint64_t v0 = UINT64_C(0x736f6d6570736575) ^ k0;
    const uint8_t *end = in + inlen - (inlen % sizeof(uint64_t));
    for (; in != end; in += 8) {
        uint64_t m = U8TO64_LE(in);
        v3 ^= m;
        SIPROUND; SIPROUND;
        v0 ^= m;
    }
    const int left = inlen & 7;
    uint64_t b = ((uint64_t)inlen) << 56;
    switch (left) {
    case 7: b |= ((uint64_t)in[6]) << 48;
    case 6: b |= ((uint64_t)in[5]) << 40;
    case 5: b |= ((uint64_t)in[4]) << 32;
    case 4: b |= ((uint64_t)in[3]) << 24;
    case 3: b |= ((uint64_t)in[2]) << 16;
    case 2: b |= ((uint64_t)in[1]) << 8;
    case 1: b |= ((uint64_t)in[0]); break;
    case 0: break;
    }
    v3 ^= b;
    SIPROUND; SIPROUND;
    v0 ^= b;
    v2 ^= 0xff;
    SIPROUND; SIPROUND; SIPROUND; SIPROUND;
    b = v0 ^ v1 ^ v2 ^ v3;
    uint64_t out = 0;
    U64TO8_LE((uint8_t*)&out, b);
    return out;
}

//-----------------------------------------------------------------------------
// MurmurHash3 was written by Austin Appleby, and is placed in the public
// domain. The author hereby disclaims copyright to this source code.
//
// Murmur3_86_128
//-----------------------------------------------------------------------------
static void MM86128(const void *key, const int len, uint32_t seed, void *out) {
#define	ROTL32(x, r) ((x << r) | (x >> (32 - r)))
#define FMIX32(h) h^=h>>16; h*=0x85ebca6b; h^=h>>13; h*=0xc2b2ae35; h^=h>>16;
    const uint8_t * data = (const uint8_t*)key;
    const int nblocks = len / 16;
    uint32_t h1 = seed;
    uint32_t h2 = seed;
    uint32_t h3 = seed;
    uint32_t h4 = seed;
    uint32_t c1 = 0x239b961b; 
    uint32_t c2 = 0xab0e9789;
    uint32_t c3 = 0x38b34ae5; 
    uint32_t c4 = 0xa1e38b93;
    const uint32_t * blocks = (const uint32_t *)(data + nblocks*16);
    for (int i = -nblocks; i; i++) {
        uint32_t k1 = blocks[i*4+0];
        uint32_t k2 = blocks[i*4+1];
        uint32_t k3 = blocks[i*4+2];
        uint32_t k4 = blocks[i*4+3];
        k1 *= c1; k1  = ROTL32(k1,15); k1 *= c2; h1 ^= k1;
        h1 = ROTL32(h1,19); h1 += h2; h1 = h1*5+0x561ccd1b;
        k2 *= c2; k2  = ROTL32(k2,16); k2 *= c3; h2 ^= k2;
        h2 = ROTL32(h2,17); h2 += h3; h2 = h2*5+0x0bcaa747;
        k3 *= c3; k3  = ROTL32(k3,17); k3 *= c4; h3 ^= k3;
        h3 = ROTL32(h3,15); h3 += h4; h3 = h3*5+0x96cd1c35;
        k4 *= c4; k4  = ROTL32(k4,18); k4 *= c1; h4 ^= k4;
        h4 = ROTL32(h4,13); h4 += h1; h4 = h4*5+0x32ac3b17;
    }
    const uint8_t * tail = (const uint8_t*)(data + nblocks*16);
    uint32_t k1 = 0;
    uint32_t k2 = 0;
    uint32_t k3 = 0;
    uint32_t k4 = 0;
    switch(len & 15) {
    case 15: k4 ^= tail[14] << 16;
    case 14: k4 ^= tail[13] << 8;
    case 13: k4 ^= tail[12] << 0;
             k4 *= c4; k4  = ROTL32(k4,18); k4 *= c1; h4 ^= k4;
    case 12: k3 ^= tail[11] << 24;
    case 11: k3 ^= tail[10] << 16;
    case 10: k3 ^= tail[ 9] << 8;
    case  9: k3 ^= tail[ 8] << 0;
             k3 *= c3; k3  = ROTL32(k3,17); k3 *= c4; h3 ^= k3;
    case  8: k2 ^= tail[ 7] << 24;
    case  7: k2 ^= tail[ 6] << 16;
    case  6: k2 ^= tail[ 5] << 8;
    case  5: k2 ^= tail[ 4] << 0;
             k2 *= c2; k2  = ROTL32(k2,16); k2 *= c3; h2 ^= k2;
    case  4: k1 ^= tail[ 3] << 24;
    case  3: k1 ^= tail[ 2] << 16;
    case  2: k1 ^= tail[ 1] << 8;
    case  1: k1 ^= tail[ 0] << 0;
             k1 *= c1; k1  = ROTL32(k1,15); k1 *= c2; h1 ^= k1;
    };
    h1 ^= len; h2 ^= len; h3 ^= len; h4 ^= len;
    h1 += h2; h1 += h3; h1 += h4;
    h2 += h1; h3 += h1; h4 += h1;
    FMIX32(h1); FMIX32(h2); FMIX32(h3); FMIX32(h4);
    h1 += h2; h1 += h3; h1 += h4;
    h2 += h1; h3 += h1; h4 += h1;
    ((uint32_t*)out)[0] = h1;
    ((uint32_t*)out)[1] = h2;
    ((uint32_t*)out)[2] = h3;
    ((uint32_t*)out)[3] = h4;
}

// hashmap_sip returns a hash value for `data` using SipHash-2-4.
uint64_t hashmap_sip(const void *data, size_t len, 
                     uint64_t seed0, uint64_t seed1)
{
    return SIP64((uint8_t*)data, len, seed0, seed1);
}

// hashmap_murmur returns a hash value for `data` using Murmur3_86_128.
uint64_t hashmap_murmur(const void *data, size_t len, 
                        uint64_t seed0, uint64_t seed1)
{
    char out[16];
    MM86128(data, len, seed0, &out);
    return *(uint64_t*)out;
}

// End of hashmap

typedef uint64_t (*hashfn)(const void* item, uint64_t seed0, uint64_t seed1);
typedef int (*comparefn)(const void* a, const void* b, void *udata);
typedef void (*freefn)(void* item);

typedef struct map {
  struct hashmap* map;
  int refs;
  int len;
} map;

map* map_new(int elemsize, comparefn comparefn, hashfn hashfn, freefn freefn) {
  map* m = malloc(sizeof(map));
  m->map = hashmap_new(elemsize, 0, 0, 0, hashfn, comparefn, freefn, NULL);
  m->refs = 1;
  return m;
}

void map_free(map* map) {
  if (map == NULL) {
    return;
  }
  map->refs--;
  if (map->refs == 0) {
    hashmap_free(map->map);
    free(map);
  }
}

// Switch-case (adapted from https://github.com/haipome/fnv)
#define FNV_32_PRIME 0x01000193
#define FNV1_32_INIT 0x811c9dc5

uint32_t str_hash(string* str) {
  uint32_t hval = FNV1_32_INIT;
  for (int i = 0; i < str->len; i++) {
    hval ^= str->data[i];
    hval *= FNV_32_PRIME;
  }
  return hval;
}

// Any
typedef void (*anyfreefn)(void*);

typedef struct any {
  void* value;
  int typ;
  int refs;
  anyfreefn freefn;
} any;

any* any_new(void* val, int typ, anyfreefn freefn) {
  any* a = malloc(sizeof(any));
  a->refs = 1;
  a->value = val;
  a->typ = typ;
  a->freefn = freefn;
  return a;
}

void any_free(any* a) {
  if (a == NULL) {
    return;
  }
  a->refs--;
  if (a->refs == 0) {
    if (a->freefn != NULL) {
      a->freefn(a->value);
    }
    free(a);
  }
}

const char* const anytyps[];

void any_try_cast(const char* pos, any* a, int typ) {
  if (a->typ != typ) {
    fprintf(stderr, "%s: cannot cast value of type %s to %s\n", pos, anytyps[a->typ], anytyps[typ]);
    exit(1);
  }
}
