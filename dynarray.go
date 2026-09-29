// Package dynarray contains data structures for Dynamic Arrays.
package dynarray

import (
	"errors"
)

const kSkipFactor = 5

var (
	kIndexMustBePositive           = errors.New("Index must be positive")
	kInvalidIndexForInsertOrDelete = errors.New(
		"Index invalid for insert or delete")
	kInsertOrDeleteCountMustBeNonNegative = errors.New(
		"Insert or delete count must be non-negative")
	kWouldOverflowLength = errors.New(
		"Insert would overflow length")
)

// DynArray works like regular arrays except for the following differences:
//
//  1. Indexing is one based.
//  2. An index can hold no value at all.
//  3. Memory usage depends on count of values stored not the index count.
//  4. When storing values beyond the highest index, length grows automatically instead of seeing index out of bounds.
//  5. Inserting into or deleting from the beginning or middle of the array takes O(log n) time instead of O(n) time.
//  6. Accessing a random array element takes O(log n) time instead of O(1) time.
//
// Cursor instances can speed up array access if reading and writing to the
// same area of the DynArray. DynArray instances are not safe to use with
// multiple goroutines.
//
// The zero value of a DynArray[T] has 0 length and is ready to use.
// DynArray instances do not support copy by value.
type DynArray[T any] struct {

	// The nodes. index 0 is a linked list of the values in this DynArray.
	// 1 is a linked list of some pointers into the linked list at index 0.
	// 2 is a linked list of some pointers into the linked list at index 1 etc.
	nodeLists []*node[T]

	// The active cursors
	cursors map[*Cursor[T]]struct{}

	ctrlNo int64

	// The length of the DynArray
	length int
}

// New creates a new DynArray of specific length with all the indexes
// containing no value. The range of indexes of the returned DynArray is
// between 1 and length inclusive. New panics if length < 0 or if length >
// math.Maxint - 1.
func New[T any](length int) *DynArray[T] {
	if length < 0 {
		panic("length must be non-negative")
	}
	if length+1 < 0 {
		panic(kWouldOverflowLength)
	}
	return &DynArray[T]{
		nodeLists: adjustNodeListsLength(([]*node[T])(nil), maxLevels(length)),
		cursors:   make(map[*Cursor[T]]struct{}),
		length:    length,
	}
}

// Len returns the length of this DynArray, which is the highest index. Len
// is not to be confused with the number of values stored in this DynArray.
func (d *DynArray[T]) Len() int {
	return d.length
}

// Truncate makes the length of this DynArray be newLength. If newLength is
// less than d.Len() any values past newLength get deleted. Truncate returns
// an error if newLength is negative or bigger than math.MaxInt - 1.
func (d *DynArray[T]) Truncate(newLength int) error {
	oldLength := d.length
	if newLength < oldLength {
		return d.Delete(newLength+1, oldLength-newLength)
	}
	return d.Insert(oldLength+1, newLength-oldLength)
}

// Insert inserts count new indexes into this DynArray right before the given
// index. The new indexes hold no value. This operation takes O(log n) time.
// Insert returns an error if index <= 0 or index > d.Len() + 1.
// index = d.Len() + 1 means insert past the end of the DynArray.
func (d *DynArray[T]) Insert(index, count int) error {
	if index <= 0 {
		return kIndexMustBePositive
	}
	c := d.NewCursorAt(index)
	defer c.Close()
	return c.Insert(count)
}

// Delete deletes count indexes from this DynArray starting at index. This
// operation takes O(log n + k) where k is the number of values deleted.
// Delete returns an error if index <= 0 or if deleting count indexes would
// go past the end of the DynArray.
func (d *DynArray[T]) Delete(index, count int) error {
	if index <= 0 {
		return kIndexMustBePositive
	}
	c := d.NewCursorAt(index)
	defer c.Close()
	return c.Delete(count)
}

// Set sets the given index of this DynArray to value. This operation
// takes O(log n) time. Set panics if index <= 0. If index points past the
// end of the DynArray then d.Len() increases to index when Set returns.
// All the indexes between the old and new length have no value.
func (d *DynArray[T]) Set(index int, value T) {
	c := d.NewCursorAt(index)
	defer c.Close()
	c.Set(value)
}

// Clear clears the given index of this DynArray so that it holds no value.
// The operation takes O(log n) time. Clear panics if index <= 0. If index
// points past the end of the DynArray then Clear is a no-op.
func (d *DynArray[T]) Clear(index int) {
	c := d.NewCursorAt(index)
	defer c.Close()
	c.Clear()
}

// Get gets the value at the given index of this DynArray. This operation
// takes O(log n) time. If the index has no value, Get returns false.
// Get panics if index <= 0. If index points past the end of the DynArray
// then Get returns false for no value.
func (d *DynArray[T]) Get(index int) (T, bool) {
	c := d.NewCursorAt(index)
	defer c.Close()
	return c.Get()
}

// GetValue works like Get except it returns just a value. GetValue returns
// the zero value if the index either has no value or the index has the zero
// value.
func (d *DynArray[T]) GetValue(index int) T {
	result, _ := d.Get(index)
	return result
}

// GetOrElse works like Get, but it returns defaultValue if index has no
// value.
func (d *DynArray[T]) GetOrElse(index int, defaultValue T) T {
	result, ok := d.Get(index)
	if ok {
		return result
	}
	return defaultValue
}

// NewCursor returns a new cursor pointing at index 1 of this DynArray.
// It is the same as d.NewCursorAt(1).
func (d *DynArray[T]) NewCursor() *Cursor[T] {
	return d.NewCursorAt(1)
}

// NewCursorAt returns a new cursor pointing at the given index of this
// DynArray. NewCursorAt panics if index <= 0.
func (d *DynArray[T]) NewCursorAt(index int) *Cursor[T] {
	if index <= 0 {
		panic(kIndexMustBePositive)
	}
	result := &Cursor[T]{
		array: d,
		index: index,
	}
	if d.cursors == nil {
		d.cursors = make(map[*Cursor[T]]struct{})
	}
	d.cursors[result] = struct{}{}
	return result
}

// FirstValue returns a cursor pointing to the first index of the this
// DynArray that has a value. If this DynArray has no actual values in it,
// FirstValue returns false.
func (d *DynArray[T]) FirstValue() (*Cursor[T], bool) {
	if len(d.nodeLists) == 0 || d.nodeLists[0].next == nil {
		return nil, false
	}
	firstIndex := d.nodeLists[0].next.indexIncr
	result := d.NewCursorAt(firstIndex)
	return result, true
}

// levelCountForNewValueNode returns how many levels up the new value node
// should extend into the node tree heirarchy. 1 means just the value node in
// linked list 0; 2 means also add a node in linked list 1 etc.
// As long as d.nodeLists has non-zero length, levelCountForNewValueNode() will
// not return a value bigger than len(d.nodeLists).
func (d *DynArray[T]) levelCountForNewValueNode() int {
	d.ctrlNo++
	result := 1
	x := d.ctrlNo
	for result < len(d.nodeLists) && x%kSkipFactor == 0 {
		result++
		x /= kSkipFactor
	}
	return result
}

func (d *DynArray[T]) setSize(newSize int) error {
	if newSize < 0 || newSize+1 < 0 {
		return kWouldOverflowLength
	}
	d.length = newSize
	d.nodeLists = adjustNodeListsLength(d.nodeLists, maxLevels(d.length))
	return nil
}

// Cursor points to a particular index in a DynArray. Accessing elements
// of the DynArray through a Cursor takes O(1) time. However creating a
// Cursor instance and using it for the first time takes O(log n) time.
// Incrementing or decrementing a cursor's index by just one takes O(1) time,
// but moving a cursor's index by a lot takes closer to O(log n) time. A
// Cursor instance can point past the end of its DynArray. Calling Set on
// such a Cursor automatically extends the length of the DynArray to include
// that new value. Indexes between the old an new length have no value.
// Cursor instances must be closed when no longer in use to avoid leaks.
type Cursor[T any] struct {
	cursorPieces []cursorPiece[T]
	array        *DynArray[T]
	index        int
	valid        bool
}

// Index returns the index in the DynArray this cursor points to.
func (c *Cursor[T]) Index() int {
	c.checkBasic()
	return c.index
}

// Len returns the length of the DynArray this cursor belongs to.
func (c *Cursor[T]) Len() int {
	c.checkBasic()
	return c.array.Len()
}

// Insert inserts count new indexes just before this cursor's index. The new
// indexes contain no value. This operation takes O(log n) time. Insert moves
// this cursor forward by count indexes making room for the indexes inserted.
// Other Cursors that point to an index at or after this cursor get moved
// forward by count indexes automatically. Insert returns an error if this
// cursor's index is more than c.Len() + 1 or if count is negative or count so
// big that it overflows the length of its DynArray.
func (c *Cursor[T]) Insert(count int) error {
	c.checkBasic()
	if count < 0 {
		return kInsertOrDeleteCountMustBeNonNegative
	}
	if c.index > c.array.Len()+1 {
		return kInvalidIndexForInsertOrDelete
	}
	if err := c.array.setSize(c.array.Len() + count); err != nil {
		return err
	}
	oldIndex := c.index
	c.resizeCursorPieces()
	c.insertOrDeleteJustBefore(count)
	c.invalidateOtherCursorsForInsert(oldIndex, count)
	return nil
}

// Delete deletes count indexes starting at this cursor's index. This
// operation takes O(k + log n) time where k is the number of actual values
// deleted. Other Cursors that point to an index after the ones that got
// deleted get moved up by count automatically. If a cursor points to an
// index that got deleted, it is automatically closed. Delete returns an
// error if there are fewer than count indexes before the end.
func (c *Cursor[T]) Delete(count int) error {
	c.checkBasic()
	if count < 0 {
		return kInsertOrDeleteCountMustBeNonNegative
	}
	if count > c.array.Len()+1-c.index {
		return kInvalidIndexForInsertOrDelete
	}
	origIndex := c.index
	for c.index < origIndex+count {
		if c.hasValue() {
			c.deleteValue()
		}
		newIndex := origIndex + count
		if c.hasNextValue() {
			newIndex = min(c.cursorPieces[0].nextIndex(), newIndex)
		}
		c.forwardToIndex(newIndex)
		c.index = newIndex
	}
	c.insertOrDeleteJustBefore(-count)
	c.array.setSize(c.array.Len() - count)
	c.resizeCursorPieces()
	c.invalidateOtherCursorsForDelete(origIndex, count)
	return nil
}

// Set sets the value at this cursor's index. This operation takes O(1) time.
// If this cursor points past the end of its DynArray, then Set also increases
// the length of this DynArray to include the value that was set.
func (c *Cursor[T]) Set(value T) {
	c.checkBasic()
	if c.index > c.array.Len() {
		if err := c.array.setSize(c.index); err != nil {
			panic(err)
		}
		c.resizeCursorPieces()
	}
	if !c.hasValue() {
		c.createValue()
	}
	c.cursorPieces[0].ptr.value = value
	c.invalidateOtherCursors()
}

// Clear clears the value at this cursor's index so that it contains no value.
// This operation takes O(1) time. If this cursor is pointing past the end of
// its DynArray, then Clear is a no-op.
func (c *Cursor[T]) Clear() {
	c.checkBasic()
	if !c.hasValue() {
		return
	}
	c.deleteValue()
	c.invalidateOtherCursors()
}

// Get gets the value at this cursor's index. This operation takes O(1) time.
// If the index has no value, Get returns false. If this cursor is pointing
// past the end of its DynArray, then Get returns false for no value.
func (c *Cursor[T]) Get() (T, bool) {
	c.checkBasic()
	if !c.hasValue() {
		var result T
		return result, false
	}
	return c.cursorPieces[0].ptr.value, true
}

// GetValue works like Get except it returns just the value. GetValue returns
// the zero value if either the index the cursor points to has no value
// or it has the zero value.
func (c *Cursor[T]) GetValue() T {
	result, _ := c.Get()
	return result
}

// GetOrElse works like Get, but it returns defaultValue if the index
// the cursor points to has no value.
func (c *Cursor[T]) GetOrElse(defaultValue T) T {
	result, ok := c.Get()
	if ok {
		return result
	}
	return defaultValue
}

// NextValue advances this curosr's index to the next actual value. If there
// is no next value, then NextValue returns false. NextValue takes O(1) time.
func (c *Cursor[T]) NextValue() bool {
	c.checkBasic()
	if !c.hasNextValue() {
		return false
	}
	newIndex := c.cursorPieces[0].nextIndex()
	c.forwardToIndex(newIndex)
	c.index = newIndex
	return true
}

// Advance advances this cursor's index forward by count. If count is negative,
// Advance moves the cursor backwards. Moving the cursor forward or backward by
// k, takes O(log k) time. Advance returns an error if the cursor's new
// index would be <= 0 or if it would overflow.
func (c *Cursor[T]) Advance(count int) error {
	c.checkBasic()
	if c.index+count <= 0 {
		return kIndexMustBePositive
	}
	if count > 0 {
		c.forwardToIndex(c.index + count)
		c.index += count
	} else if count < 0 {
		c.backwardToIndex(c.index + count)
		c.index += count
	}
	return nil
}

// Goto changes the index of this cursor to be index. Goto delegates to
// Advance. Goto panics if index <= 0.
func (c *Cursor[T]) Goto(index int) {
	if err := c.Advance(index - c.index); err != nil {
		panic(err)
	}
}

// Close closes this cursor. Close can be called multiple times on the same
// Cursor. Calling methods on a closed cursor besides the Close and IsClosed
// methods panics.
func (c *Cursor[T]) Close() {
	delete(c.array.cursors, c)
}

// IsClosed returns true if this cursor has been closed
func (c *Cursor[T]) IsClosed() bool {
	_, ok := c.array.cursors[c]
	return !ok
}

func (c *Cursor[T]) hasNextValue() bool {
	return len(c.cursorPieces) > 0 && c.cursorPieces[0].ptr.next != nil
}

func (c *Cursor[T]) invalidateOtherCursors() {
	for cursor := range c.array.cursors {
		if cursor == c {
			continue
		}
		cursor.valid = false
	}
}

func (c *Cursor[T]) invalidateOtherCursorsForInsert(insertPoint, count int) {
	for cursor := range c.array.cursors {
		if cursor == c {
			continue
		}
		cursor.valid = false
		if cursor.index >= insertPoint {
			cursor.index += count
		}
	}
}

func (c *Cursor[T]) invalidateOtherCursorsForDelete(deletePoint, count int) {
	for cursor := range c.array.cursors {
		if cursor == c {
			continue
		}
		if cursor.index >= deletePoint && cursor.index < deletePoint+count {
			delete(c.array.cursors, cursor)
			continue
		}
		cursor.valid = false
		if cursor.index >= deletePoint {
			cursor.index -= count
		}
	}
}

func (c *Cursor[T]) hasValue() bool {
	return len(c.cursorPieces) != 0 && c.index == c.cursorPieces[0].indexSoFar
}

func (c *Cursor[T]) createValue() {
	levelCount := c.array.levelCountForNewValueNode()
	for idx := 0; idx < levelCount; idx++ {
		c.cursorPieces[idx].newNode(c.index)
		if idx > 0 {
			c.cursorPieces[idx].ptr.down = c.cursorPieces[idx-1].ptr
		}
	}
}

func (c *Cursor[T]) deleteValue() {
	c.backwardToIndex(c.index - 1)
	for idx := range c.cursorPieces {
		cp := &c.cursorPieces[idx]
		if cp.ptr.next == nil || cp.nextIndex() != c.index {
			break
		}
		cp.removeNextNode()
	}
}

func (c *Cursor[T]) insertOrDeleteJustBefore(count int) {
	for idx := range c.cursorPieces {
		cp := &c.cursorPieces[idx]
		if c.index == cp.indexSoFar {
			cp.indexSoFar += count
			cp.ptr.indexIncr += count
		} else if cp.ptr.next != nil {
			cp.ptr.next.indexIncr += count
		}
	}
	c.index += count
}

func (c *Cursor[T]) resetCursorPieces() {
	c.updateCursorPieces(true)
}

func (c *Cursor[T]) resizeCursorPieces() {
	c.updateCursorPieces(false)
}

func (c *Cursor[T]) resetCursorPiece(idx int) {
	c.cursorPieces[idx].ptr = c.array.nodeLists[idx]
	c.cursorPieces[idx].indexSoFar = 0
}

func (c *Cursor[T]) updateCursorPieces(fullReset bool) {
	if fullReset {
		for idx := range min(len(c.array.nodeLists), len(c.cursorPieces)) {
			c.resetCursorPiece(idx)
		}
	}
	for len(c.array.nodeLists) > len(c.cursorPieces) {
		ln := len(c.cursorPieces)
		c.cursorPieces = append(
			c.cursorPieces, cursorPiece[T]{ptr: c.array.nodeLists[ln]})
	}
	for len(c.array.nodeLists) < len(c.cursorPieces) {
		ln := len(c.cursorPieces)
		c.cursorPieces[ln-1].ptr = nil
		c.cursorPieces = c.cursorPieces[:ln-1]
	}
}

func (c *Cursor[T]) checkBasic() {
	c.checkNotClosed()
	c.checkCursorPieces()
}

func (c *Cursor[T]) checkCursorPieces() {
	if !c.valid {
		c.resetCursorPieces()
		c.forwardToIndex(c.index)
		c.valid = true
	}
}

func (c *Cursor[T]) checkNotClosed() {
	if c.IsClosed() {
		panic("Cursor closed")
	}
}

func (c *Cursor[T]) findFirstNoForwardLevel(desiredIndex int) int {
	var idx int
	for idx = 0; idx < len(c.cursorPieces); idx++ {
		if c.cursorPieces[idx].ptr.next == nil || c.cursorPieces[idx].nextIndex() > desiredIndex {
			break
		}
	}
	return idx
}

func (c *Cursor[T]) findFirstNoBackwardLevel(desiredIndex int) int {
	var idx int
	for idx = 0; idx < len(c.cursorPieces); idx++ {
		if c.cursorPieces[idx].indexSoFar <= desiredIndex {
			break
		}
	}
	return idx
}

func (c *Cursor[T]) forwardToIndex(desiredIndex int) {
	firstNoForwardLevel := c.findFirstNoForwardLevel(desiredIndex)
	for idx := firstNoForwardLevel - 1; idx >= 0; idx-- {
		c.cursorPieces[idx].forwardToIndex(desiredIndex)
		if idx > 0 {
			c.cursorPieces[idx-1].ptr = c.cursorPieces[idx].ptr.down
			c.cursorPieces[idx-1].indexSoFar = c.cursorPieces[idx].indexSoFar
		}
	}
}

func (c *Cursor[T]) backwardToIndex(desiredIndex int) {
	firstNoBackwardLevel := c.findFirstNoBackwardLevel(desiredIndex)
	for idx := firstNoBackwardLevel - 1; idx >= 0; idx-- {
		if idx+1 == len(c.cursorPieces) || c.cursorPieces[idx+1].isAnchor() {
			c.resetCursorPiece(idx)
		} else {
			c.cursorPieces[idx].ptr = c.cursorPieces[idx+1].ptr.down
			c.cursorPieces[idx].indexSoFar = c.cursorPieces[idx+1].indexSoFar
		}
		c.cursorPieces[idx].forwardToIndex(desiredIndex)
	}
}

type node[T any] struct {
	down      *node[T]
	next      *node[T]
	indexIncr int
	value     T
}

type cursorPiece[T any] struct {
	ptr        *node[T]
	indexSoFar int
}

func (c *cursorPiece[T]) isAnchor() bool {
	return c.indexSoFar == 0
}

func (c *cursorPiece[T]) forwardToIndex(desiredIndex int) {
	for c.ptr.next != nil && c.nextIndex() <= desiredIndex {
		c.ptr = c.ptr.next
		c.indexSoFar += c.ptr.indexIncr
	}
}

func (c *cursorPiece[T]) nextIndex() int {
	return c.indexSoFar + c.ptr.next.indexIncr
}

func (c *cursorPiece[T]) newNode(index int) {
	afterNode := c.ptr.next
	diff := index - c.indexSoFar
	newnode := &node[T]{indexIncr: diff}
	c.ptr.next = newnode
	if afterNode != nil {
		newnode.next = afterNode
		afterNode.indexIncr -= diff
	}
	c.ptr = newnode
	c.indexSoFar = index
}

func (c *cursorPiece[T]) removeNextNode() {
	afterAfterNode := c.ptr.next.next
	if afterAfterNode != nil {
		afterAfterNode.indexIncr += c.ptr.next.indexIncr
	}
	c.ptr.next.next = nil
	c.ptr.next.down = nil
	c.ptr.next = afterAfterNode
}

// maxLevels returns the maximum number of levels of linked lists given the
// length of the DynArray.
func maxLevels(length int) int {
	result := 0
	for length > 0 {
		result++
		length /= kSkipFactor
	}
	return result
}

func adjustNodeListsLength[T any](nodeLists []*node[T], length int) []*node[T] {
	for length > len(nodeLists) {
		nodeLists = append(nodeLists, &node[T]{})
	}
	for length < len(nodeLists) {
		ln := len(nodeLists)
		nodeLists[ln-1] = nil
		nodeLists = nodeLists[:ln-1]
	}
	return nodeLists
}
