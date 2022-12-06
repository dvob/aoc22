#include <unistd.h>
#include <fcntl.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <errno.h>
#include <stdbool.h>

char *read_all(int fd) {
	char *data = NULL;
	char *new_data = NULL;
	char buf[4096];
	int len = 0;
	for (;;) {
		int bytes_read = read(fd, buf, 4096);
		if ( bytes_read == -1 ) {
			goto error;
		}
		if ( bytes_read == 0 ) {
			break;
		}
		new_data = reallocarray(data, sizeof(char), len + bytes_read);
		if ( new_data == NULL ) {
			goto error;
		}
		data = new_data;
		memcpy(data + len, buf, bytes_read);
		len += bytes_read;
	}
	if ( data[len-1] != '\n' ) {
		errno = EINVAL;
		goto error;
	}
	data[len-1] = '\0';
	return data;
error:
	if ( data != NULL ) {
		free(data);
	}
	return NULL;
}

int solve(char *input, int num) {
	char *last = calloc(sizeof(char), num);
	if ( last == NULL ) {
		perror("no memory");
		return -1;
	}
	int pos = 0;
	for ( ; *input; input++ ) {
		last[pos % num] = *input;
		pos++;
		if ( pos < num + 1 ) {
			continue;
		}

		bool found_duplicate = false;
		for ( int i = 0; i < num; i++ ) {
			for ( int j = 0; j < num; j++ ) {
				if ( j == i ) {
					continue;
				}

				if ( last[i] == last[j] ) {
					found_duplicate = true;
					break;
				}
			}
			if ( found_duplicate ) {
				break;
			}
		}
		if ( !found_duplicate ) {
			free(last);
			return pos;
		}
	}
	free(last);
	return -1;
}

int main(int argc, char **argv) {
	if ( argc < 2 ) {
		fprintf(stderr, "missing argument: filename\n");
		return EXIT_FAILURE;
	}

	int fd = open(argv[1], O_RDONLY);
	if ( fd == -1 ) {
		perror("failed to open file");
		return EXIT_FAILURE;
	}

	char *input = read_all(fd);
	if ( input == NULL ) {
		perror("failed to read file");
		return EXIT_FAILURE;
	}

	int result1 = solve(input, 4);
	printf("%d\n", result1);

	int result2 = solve(input, 14);
	printf("%d\n", result2);

	free(input);
}
