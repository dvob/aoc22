#include <stdlib.h>
#include <fcntl.h>
#include <unistd.h>
#include <stdbool.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <ctype.h>

#include "list.h"

typedef struct inst {
	int qty;
	int src;
	int dst;
} inst;

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

int parse_int(slice s, int *i) {
	char tmp[20];
	if ( s.ptr == NULL ) {
		return EXIT_FAILURE;
	}
	if ( s.len > 19 ) {
		return EXIT_FAILURE;
	}
	memcpy(tmp, s.ptr, s.len);
	tmp[s.len] = '\0';
	long l = strtol(tmp, NULL, 10);
	if ( errno != 0 ) {
		return EXIT_FAILURE;
	}
	*i = (int)l;
	return EXIT_SUCCESS;
}

list *get_lines(char *s, int len) {
	list *l = list_new(sizeof(slice));
	iterator it = iter_new(s, len, '\n');

	slice line;
	for (;;) {
		line = iter_next(&it);
		if ( line.ptr == NULL ) {
			break;
		}
		int ret = list_push(l, (void *)&line);
		if ( ret != EXIT_SUCCESS ) {
			list_free(l);
			return NULL;
		}
	}
	return l;
}

int find_empty_line(list *lines) {
	for (int i = 0; i < list_len(lines); i++ ) {
		slice s = *(slice *)list_get(lines, i);
		if ( s.len == 0 ) {
			return i;
		}
	}
	return -1;
}

void print_stacks(list *l) {
	// list stacks
	for (int i = 0; i < list_len(l); i++) {
		void *ptr = list_get(l, i);
		list *stack = *(list **)list_get(l, i);
		printf("stack %d: ", i+1);
		for (int j = 0; j < list_len(stack); j++ ) {
			char c = *(char *)list_get(stack, j);
			printf("%c ", c);
		}
		printf("\n");
	}
}

void print_solution(list *l) {
	for (int i = 0; i < list_len(l); i++) {
		void *ptr = list_get(l, i);
		list *stack = *(list **)list_get(l, i);
		if ( list_len(stack) > 0 ) {
			char c = *(char *)list_get(stack, list_len(stack)-1);
			printf("%c", c);
		}
	}
	printf("\n");
}

list *get_stacks(list *lines) {
	// get offsets
	int head_end = find_empty_line(lines);
	if ( head_end == -1 ) {
		fprintf(stderr, "no empty line found\n");
		return NULL;
	}
	slice last_stack_line = *(slice *)list_get(lines, head_end - 1);

	list *offsets = list_new(sizeof(int));
	if ( offsets == NULL ) {
		perror("no memory");
		return NULL;
	}

	for (int i = 0; i < last_stack_line.len; i++) {
		if ( isspace(last_stack_line.ptr[i]) ) {
			continue;
		}

		int ret = list_push(offsets, (void *)&i);
		if ( ret != EXIT_SUCCESS ) {
			list_free(offsets);
			return NULL;
		}

		int j;
		for (j = i; j < last_stack_line.len; j++) {
			if ( isspace(last_stack_line.ptr[j]) ) {
				break;
			}
		}
		i = j;
	}

	// initialize stacks
	list *stacks = list_new(sizeof(list *));
	stacks->free = list_free_func;

	for (int j = 0; j < list_len(offsets); j++) {
		list *l = list_new(sizeof(char));
		if (l == NULL) {
			list_free(stacks);
			return NULL;
		}
		int ret = list_push(stacks, (void *)&l);
		if ( ret != EXIT_SUCCESS ) {
			list_free(stacks);
			return NULL;
		}
	}

	// fill stacks from bottom
	for ( int i = head_end - 2; i >= 0; i-- ) {
		// take each offset
		slice line = *(slice *)list_get(lines, i);
		for (int j = 0; j < list_len(offsets); j++) {
			int offset = list_get_int(offsets, j);
			if ( !isspace(line.ptr[offset]) ) {
				list *stack = *(list **)list_get(stacks, j);
				int ret = list_push(stack, (void *)&line.ptr[offset]);
				if ( ret != EXIT_SUCCESS ) {
					return NULL;
				}
			}
		}
	}
	list_free(offsets);
	return stacks;
}

list *get_instructions(list *lines) {
	int inst_start = find_empty_line(lines);
	if ( inst_start == -1 ) {
		fprintf(stderr, "no empty line found\n");
		return NULL;
	}
	inst_start++;

	list *insts = list_new(sizeof(inst));
	for (int i = inst_start; i < list_len(lines); i++) {
		slice line = *(slice *)list_get(lines, i);
		inst in = {};
		// read instruction
		int ret = sscanf(line.ptr, "move %d from %d to %d", &in.qty, &in.src, &in.dst);
		if ( ret != 3 ) {
			list_free(insts);
			return NULL;
		}
		ret = list_push(insts, (void *)&in);
		if ( ret != EXIT_SUCCESS ) {
			list_free(insts);
			return NULL;
		}
	}
	return insts;
}

void run_instructions(list *insts, list *stacks) {
	for (int i = 0; i < list_len(insts); i++) {
		inst in = *(inst *)list_get(insts, i);
		list *src = *(list **)list_get(stacks, in.src - 1);
		list *dst = *(list **)list_get(stacks, in.dst - 1);
		for (int j = 0; j < in.qty; j++) {
			char c = *(char *)list_pop(src);
			list_push(dst, (void *)&c);
		}
	}
}

void run_instructions2(list *insts, list *stacks) {
	for (int i = 0; i < list_len(insts); i++) {
		inst in = *(inst *)list_get(insts, i);
		list *src = *(list **)list_get(stacks, in.src - 1);
		list *dst = *(list **)list_get(stacks, in.dst - 1);

		// copy items to dest
		char *c = (char *)list_get(src, src->len - in.qty);
		list_extend(dst, (void *)c, in.qty);

		// remove items from source
		for (int j = 0; j < in.qty; j++) {
			list_pop(src);
		}
	}
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

	list *input = read_all(fd);
	if ( input == NULL ) {
		perror("failed to read file");
		return EXIT_FAILURE;
	}

	//printf("input='%.*s'", input->len, (char *)input->data);

	list *lines = get_lines((char *)input->data, input->len);
	if ( lines == NULL ) {
		perror("no memory");
		list_free(input);
		return EXIT_FAILURE;
	}

	list *stacks = get_stacks(lines);
	if ( stacks == NULL ) {
		fprintf(stderr, "failed to get stacks\n");
		goto error;
	}

	//print_stacks(stacks);

	list *insts = get_instructions(lines);
	if ( insts == NULL ) {
		goto error;
	}

	run_instructions(insts, stacks);
	print_solution(stacks);
	list_free(stacks);

	stacks = get_stacks(lines);
	if ( stacks == NULL ) {
		goto error;
	}
	run_instructions2(insts, stacks);
	print_solution(stacks);
	list_free(stacks);

	list_free(insts);
	list_free(lines);
	list_free(input);
	return EXIT_SUCCESS;

error:
	if ( input != NULL ) {
		list_free(input);
	}
	if ( lines != NULL ) {
		list_free(lines);
	}
	if ( insts != NULL ) {
		list_free(insts);
	}
	if ( stacks != NULL ) {
		list_free(stacks);
	}
	return EXIT_FAILURE;
}
