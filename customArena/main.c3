module arena_allocator;

import std::io;

const DEFAULT_ALIGNMENT = 2 * (void*::size);

struct Arena {
	void* base;
	usz size;
	usz offset;
	usz commited;
}

fn void* arena_alloc_aligned(Arena* a, usz size, usz alignment) {
	uptr curr_ptr = (uptr)a.base + (uptr)a.offset;
	uptr offset   = align_forward(curr_ptr, alignment);
	offset -= (uptr)a.base;

	if (offset + size > a.size) return (void*)0;

	a.commited += size;
	void* ptr = (ichar*)a.base + offset;
	a.offset = offset + size;

	return ptr;
}

macro bool is_power_of_two(x) {
	return x != 0 && (x & (x - 1)) == 0;
}

fn uptr align_forward(uptr ptr, usz alignment) {
	uptr p, a, modulo;
	if (!is_power_of_two(alignment)) {
		return 0;
	}
	p = ptr;
	a = (uptr)alignment;
	modulo = p & (a - 1);

	if (modulo) {
		p += a - modulo;
	}
	return p;
}

fn void* arena_alloc(usz size, void* context) {
	if (!size) return (void*)0;
	return arena_alloc_aligned((Arena*)context, size, DEFAULT_ALIGNMENT);
}

fn void arena_free_all(void* context) {
	Arena* a = context;
	a.offset = 0;
	a.commited = 0;
}

fn Arena arena_init(void* buffer, usz size) => (Arena) { .base = buffer, .size = size };
