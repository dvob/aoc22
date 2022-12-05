#include <stdio.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

typedef struct list {
	void *data;
	int len;
	int cap;
	unsigned long type_size;
        void (*free)(void *ptr);
} list;

list *list_new(unsigned long type_size) {
	list *l = (list *)malloc(sizeof(list));
	if ( l == NULL ) {
		return NULL;
	}
	l->data = NULL;
	l->len = 0;
	l->cap = 0;
	l->type_size = type_size;
	l->free = NULL;
	return l;
}

int list_len(list *l) {
	return l->len;
}

void *list_get(list *l, int index) {
	if ( index >= l->len ) {
		return NULL;
	}
	return l->data + l->type_size * index;
}


int list_get_int(list *l, int index) {
	int *i = (int *)list_get(l, index);
	if ( i == NULL ) {
		return 0;
	}
	return *i;
}


void list_free_item(list *l, int index) {
	if ( l->free != NULL ) {
		void *item = list_get(l, index);
		l->free(item);
	}
}

void list_free(list *l) {
	// free individual elements if a free function is set
	if ( l->free != NULL ) {
		int len = list_len(l);
		for ( int i = 0; i < len; i++ ) {
			list_free_item(l, i);
		}
	}
	free(l->data);
	free(l);
}

void list_empty(list *l) {
	for ( int i = 0; i < l->len; i++ ) {
		list_free_item(l, i);
	}
	l->len = 0;
}

void list_free_func(void *ptr) {
	list *l = *(list **)ptr;
	list_free(l);
}

int list_push(list *l, void *item) {
	if ( l->cap <= l->len ) {
		int newCap;
		if ( l->cap == 0 ) {
			newCap = 4;
		} else {
			newCap = l->cap * 2;
		}
		l->cap = newCap;
		void *new_data = reallocarray(l->data, l->type_size, l->cap); 
		if ( new_data == NULL ) {
			return EXIT_FAILURE;
		}
		l->data = new_data;
	}
	void *dst = l->data + l->type_size * l->len;
	memcpy(dst, item, l->type_size);
	l->len++;
	return EXIT_SUCCESS;
}

int list_extend(list *l, void *ptr, int len) {
	int added_items = 0;
	for (int i = 0; i < len; i++) {
		int ret = list_push(l, ptr);
		if ( ret != EXIT_SUCCESS ) {
			return added_items;
		}
		ptr += l->type_size;
		added_items++;
	}
	return added_items;
}

int list_get_char(list *l, int index) {
	char *c = (char *)list_get(l, index);
	if ( c == NULL ) {
		return 0;
	}
	return *c;
}

void *list_pop(list *l) {
	if ( l->len < 1 ) {
		return NULL;
	}
	int index = l->len - 1;
	void *item = list_get(l, index);
	l->len--;
	list_free_item(l, index);
	return item;
}

list *read_all(int fd) {
	list *l = list_new(sizeof(char));
	if ( l == NULL ) {
		return NULL;
	}

	char buf[4096];
	for (;;) {
		int ret = read(fd, &buf, 4096);
		if ( ret == -1 ) {
			goto error;
		}
		if ( ret == 0 ) {
			break;
		}
		int len = list_extend(l, buf, ret);
		if ( len != ret ) {
			goto error;
		}
	}
	return l;
error:
	list_free(l);
	return NULL;
}

