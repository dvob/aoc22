
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <stdbool.h>

typedef struct list {
	char *data;
	size_t cap;
	size_t len;
} list;

int list_add(list *l, char *buf, size_t len) {
	size_t requiredSize, newSize;
	char *newData;
	requiredSize = len + l->len;
	if ( requiredSize > l->cap ) {
		if ( l->cap == 0 ) {
			l->cap = 4;
		}
		newSize = l->cap * 2;
		while ( newSize < requiredSize ) {
			newSize *= 2;
		}
		// grow to newSize
		newData = malloc(newSize);
		if ( newData == NULL ) {
			return EXIT_FAILURE;
		}
		memcpy((void *) newData, l->data, l->len);
		free(l->data);
		l->data = newData;
		l->cap = newSize;
	}

	memcpy(l->data + l->len, buf, len);
	l->len += len;
	return EXIT_SUCCESS;
}

void list_reset(list *l) {
	if ( l->cap > 0 ) {
		free(l->data);
	}
	l->cap = 0;
	l->len = 0;
}

int read_all(int fd, list *l) {
	size_t fileSize;
	struct stat st;
	char buf[512];
	// fstat(fd, &st);
	// fileSize = st.st_size;
	for (;;) {
		int len = read(fd, buf, 512);
		if ( len == 0 ) {
			break;
		}
		if (list_add(l, buf, len) != 0) {
			return EXIT_FAILURE;
		}
	}
	return EXIT_SUCCESS;
}

int find_newline(char *data, int len, int start) {
	for ( int i = start; i < len; i++ ) {
		if ( data[i] == '\n' ) {
			return i;
		}
	}
	return -1;
}

char solveLine(char *data, int len) {
	if ( len % 2 != 0 ) {
		fprintf(stderr, "invalid line length");
		return '\0';
	}
	int partLen = len / 2;
	for ( int i = 0; i < partLen; i++ ) {
		for ( int j = partLen; j < len; j++ ) {
			if ( data[i] == data[j] ) {
				return data[i];
			}
		}
	}
	return '\0';
}

int getPrio(char c) {
	if ( c >= 'a' && c <= 'z' ) {
		return c - 'a' + 1;
	}
	if ( c >= 'A' && c <= 'Z' ) {
		return c - 'A' + 27;
	}
	return -1;
}

int solve1(char *data, int len) {
	int i = 0;
	int prioSum = 0;
	while ( i < len ) {
		int nextNewline = find_newline(data, len, i);
		if ( nextNewline == -1 ) {
			fprintf(stderr, "no newline found after %d\n", i);
			return EXIT_FAILURE;
		}
		int lineLen = nextNewline - i;
		char mixedItem = solveLine(data + i, lineLen);
		if ( mixedItem == '\0' ) {
			fprintf(stderr, "no items mixed up after %d\n", i);
			return EXIT_FAILURE;
		}

		int prio = getPrio(mixedItem);
		if ( prio == -1 ) {
			fprintf(stderr, "could not get prio for %c\n", mixedItem);
			return EXIT_FAILURE;
		}
		prioSum += prio;
	
		i = nextNewline + 1;
	}
	return prioSum;
}

bool line_contains(char *line, char c) {
	int i = 0;
	while(line[i] != '\n') {
		if ( line[i] == c ) {
			return true;
		}
		i++;
	}
	return false;
}

void print_line(char *line) {
	for ( int i = 0; line[i] != '\n'; i++ ) {
		printf("%c", line[i]);
	}
	printf("\n");
}

char get_group_badge(char *lines[3]) {
	int i = 0;
	char *first = lines[0];
	while (first[i] != '\n') {
		if ( line_contains(lines[1], first[i]) && line_contains(lines[2], first[i]) ) {
				return first[i];
		}
		i++;
	}
	return '\0';
}

int solve2(char *data, int len) {
	int i = 0;
	int prioSum = 0;
	while ( i < len ) {
		char *lines[3];
		for ( int j = 0; j < 3; j++ ) {
			lines[j] = data + i;
			int nextNewline = find_newline(data, len, i);
			if ( nextNewline == -1 ) {
				fprintf(stderr, "no newline found after %d\n", i);
				return EXIT_FAILURE;
			}
			i = nextNewline + 1;
		}
		char badgeItem = get_group_badge(lines);
		if ( badgeItem == '\0' ) {
			fprintf(stderr, "no badge item found after %d\n", i);
			return EXIT_FAILURE;
		}

		int prio = getPrio(badgeItem);
		if ( prio == -1 ) {
			fprintf(stderr, "could not get prio for %c\n", badgeItem);
			return EXIT_FAILURE;
		}
		prioSum += prio;
	
	}
	return prioSum;
}

int main(int argc, char **argv) {
	if ( argc < 2 ) {
		fprintf(stderr, "missing argument: filename\n");
		return EXIT_FAILURE;
	}

	int fd = open(argv[1], O_RDONLY);
	if ( fd == -1 ) {
		perror("file open");
		return EXIT_FAILURE;
	}

	list l;
	memset(&l, 0, sizeof(list));

	int rt = read_all(fd, &l);
	if ( rt != 0 ) {
		perror("read failed");
		return EXIT_FAILURE;
	}

	int result = solve1(l.data, l.len);
	printf("%d\n", result);

	result = solve2(l.data, l.len);
	printf("%d\n", result);

	return EXIT_SUCCESS;
}
