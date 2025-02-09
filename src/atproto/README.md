# ATProto

This is a (very minimal) reproduction of ATProto as a Go client library.

The main goal of this project is to be able to use ATProto easily to post
some content automatically.

## Main Goals

- [x] Create / Refresh Sessions
- [x] Post a record
- [x] Automatically detect facets from rich text
  - This is partial, and doesn't currently handle tags. 
  - Mentions and links are looking great, tho.
- [ ] Create nicely nested threads
- [ ] world domination

## Organization

- `atproto/interfaces`: Interface type definitions, based on the [TS client](https://github.com/bluesky-social/atproto/blob/main/packages/api/src/client/index.ts)
- `atproto/client`: A working implementation of the interfaces using [`indigo`](https://github.com/bluesky-social/indigo) libraries. 
- `atproto/richtext`: A nice implementation of the RichText type with automatic facet detection. Shamelessly copied from the [TS implemetation](https://github.com/bluesky-social/atproto/blob/main/packages/api/src/rich-text).
- `atproto/mock`: Minimal mocked types for testing purposes.