#include <stdlib.h>
#include <fcntl.h>
#include <unistd.h>
#include <stdbool.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>

typedef struct slice {
	char *ptr;
	int len;
} slice;

typedef struct iterator {
	int index;
	char *data;
	int len;
	char sep;
} iterator;

iterator iter_new(char *data, int len, char sep) {
	iterator i = {
		.index = 0,
		.data = data,
		.len = len,
		.sep = sep,
	};
	return i;
}

slice iter_next(iterator *it) {
	slice s = {
		.ptr = NULL,
		.len = 0,
	};
	if (it->index >= it->len) {
		return s;
	}
	s.ptr = it->data + it->index;
	int i;
	for ( i = it->index; i < it->len; i++) {
		if ( it->data[i] == it->sep ) {
			break;
		}
	}
	s.len = i - it->index;
	it->index = i + 1;
	return s;
};

typedef struct range {
	int from;
	int to;
} range;

// contains range a range b
bool contains_range(range a, range b) {
	return b.from >= a.from && b.to <= a.to;
}

bool do_overlap(range a, range b) {
	bool a_start_in_b = a.from >= b.from && a.from <= b.to;
	bool a_end_in_b = a.to >= b.from && a.to <= b.to;
	bool a_contains_b = contains_range(a, b);
	return a_start_in_b || a_end_in_b || a_contains_b;
}

typedef struct pair {
	range first;
	range second;
} pair;

typedef enum result_t {
	SUCCESS = 0,
	E_STDLIB = 1,
	E_PARSE = 2,
} result_t;

void print_slice(slice s) {
	if (s.ptr == NULL) {
		printf("NULL\n");
		return;
	}
	printf("'");
	for (int i = 0; i < s.len; i++) {
		printf("%c", s.ptr[i]);
	}
	printf("'\n");
}

result_t read_all(int fd, char **data, int *len) {
	result_t err = E_STDLIB;
	char buf[4096];
	void *tmp = NULL;
	*len = 0;
	for (;;) {
		int ret = read(fd, &buf, 4096);
		if ( ret == -1 ) {
			goto error;
		}
		if ( ret == 0 ) {
			break;
		}
		void *tmp_new = realloc(tmp, *len + ret);
		if ( tmp_new == NULL ) {
			goto error;
		}
		tmp = tmp_new;
		memcpy(tmp + *len, (void *)buf, ret);
		*len += ret;
	}
	*data = (char *)tmp;
	return SUCCESS;
error:
	if ( tmp != NULL ) {
		free(tmp);
	}
	return err;
}

result_t parse_int(slice s, int *i) {
	char tmp[20];
	if ( s.ptr == NULL ) {
		return E_PARSE;
	}
	if ( s.len > 19 ) {
		return E_PARSE;
	}
	memcpy(tmp, s.ptr, s.len);
	tmp[s.len] = '\0';
	long l = strtol(tmp, NULL, 10);
	if ( errno != 0 ) {
		return E_PARSE;
	}
	*i = (int)l;
	return SUCCESS;
}

result_t parse_range(slice s, range *r) {
	result_t res;

	iterator it = iter_new(s.ptr, s.len, '-');

	s = iter_next(&it);
	if ( s.ptr == NULL ) {
		return E_PARSE;
	}
	res = parse_int(s, &(*r).from);
	if ( res != SUCCESS ) {
		return res;
	}

	s = iter_next(&it);
	if ( s.ptr == NULL ) {
		return E_PARSE;
	}
	res = parse_int(s, &(*r).to);
	return res;
}

result_t parse_pair(slice s, pair *p) {
	result_t res;

	iterator it = iter_new(s.ptr, s.len, ',');

	s = iter_next(&it);
	if ( s.ptr == NULL ) {
		return E_PARSE;
	}
	res = parse_range(s, &(*p).first);
	if ( res != SUCCESS ) {
		return res;
	}

	s = iter_next(&it);
	if ( s.ptr == NULL ) {
		return E_PARSE;
	}
	res = parse_range(s, &(*p).second);
	return res;
}

int main(int argc, char *argv[]) {
	if ( argc < 2 ) {
		fprintf(stderr, "missing argument: filename\n");
		return EXIT_FAILURE;
	}

	int fd = open(argv[1], O_RDONLY);
	if ( fd == -1 ) {
		perror("failed to open file");
		return EXIT_FAILURE;
	}

	char *data;
	int len;
	result_t res = read_all(fd, &data, &len);
	if ( res != SUCCESS ) {
		perror("failed to read file");
		return EXIT_FAILURE;
	}

	iterator it = iter_new(data, len, '\n');
	int result1 = 0;
	int result2 = 0;
	for (;;) {
		slice s = iter_next(&it);
		if ( s.ptr == NULL ) {
			break;
		}
		pair p;
		res = parse_pair(s, &p);
		if ( contains_range(p.first, p.second) || contains_range(p.second, p.first) ) {
			result1++;
		}

		if ( do_overlap(p.first, p.second) ) {
			result2++;
		}
	}
	printf("%d\n", result1);
	printf("%d\n", result2);

	free(data);
}
