// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kv/kvserver/concurrency/lock/locking.proto

package lock

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// Strength represents the different locking modes that determine how key-values
// can be accessed by concurrent transactions.
//
// Locking modes apply to locks that are held with a per-key granularity. It is
// up to users of the key-value layer to decide on which keys to acquire locks
// for when imposing structure that can span multiple keys, such as SQL rows
// (see column families and secondary indexes).
//
// Locking modes have differing levels of strength, growing from "weakest" to
// "strongest" in the order that the variants are presented in the enumeration.
// The "stronger" a locking mode, the more protection it provides for the lock
// holder but the more restrictive it is to concurrent transactions attempting
// to access the same keys.
//
// Compatibility Matrix
//
// The following matrix presents the compatibility of locking strengths with one
// another. A cell with an X means that the two strengths are incompatible with
// each other and that they can not both be held on a given key by different
// transactions, concurrently. A cell without an X means that the two strengths
// are compatible with each other and that they can be held on a given key by
// different transactions, concurrently.
//
//  +-----------+-----------+-----------+-----------+-----------+
//  |           |   None    |  Shared   |  Upgrade  | Exclusive |
//  +-----------+-----------+-----------+-----------+-----------+
//  | None      |           |           |           |     X^†   |
//  +-----------+-----------+-----------+-----------+-----------+
//  | Shared    |           |           |     X     |     X     |
//  +-----------+-----------+-----------+-----------+-----------+
//  | Upgrade   |           |     X     |     X     |     X     |
//  +-----------+-----------+-----------+-----------+-----------+
//  | Exclusive |     X^†   |     X     |     X     |     X     |
//  +-----------+-----------+-----------+-----------+-----------+
//
// [†] reads under optimistic concurrency control in CockroachDB only conflict
// with Exclusive locks if the read's timestamp is equal to or greater than the
// lock's timestamp. If the read's timestamp is below the Exclusive lock's
// timestamp then the two are compatible.
type Strength int32

const (
	// None represents the absence of a lock or the intention to acquire locks.
	// It corresponds to the behavior of transactions performing key-value reads
	// under optimistic concurrency control. No locks are acquired on the keys
	// read by these requests when they evaluate. However, the reads do respect
	// Exclusive locks already held by other transactions at timestamps equal to
	// or less than their read timestamp.
	//
	// Optimistic concurrency control (OCC) can improve performance under some
	// workloads because it avoids the need to perform any locking during reads.
	// This can increase the amount of concurrency that the system can permit
	// between ongoing transactions. However, OCC does mandate a read validation
	// phase if/when transactions need to commit at a different timestamp than
	// they performed all reads at. CockroachDB calls this a "read refresh",
	// which is implemented by the txnSpanRefresher. If a read refresh fails due
	// to new key-value writes that invalidate what was previously read,
	// transactions are forced to restart. See the comment on txnSpanRefresher
	// for more.
	None Strength = 0
	// Shared (S) locks are used by read-only operations and allow concurrent
	// transactions to read under pessimistic concurrency control. Shared locks
	// are compatible with each other but are not compatible with Upgrade or
	// Exclusive locks. This means that multiple transactions can hold a Shared
	// lock on the same key at the same time, but no other transaction can
	// modify the key at the same time. A holder of a Shared lock on a key is
	// only permitted to read the key's value while the lock is held.
	//
	// Share locks are currently unused, as all KV reads are currently performed
	// optimistically (see None).
	Shared Strength = 1
	// Upgrade (U) locks are a hybrid of Shared and Exclusive locks which are
	// used to prevent a common form of deadlock. When a transaction intends to
	// modify existing KVs, it is often the case that it reads the KVs first and
	// then attempts to modify them. Under pessimistic concurrency control, this
	// would correspond to first acquiring a Shared lock on the keys and then
	// converting the lock to an Exclusive lock when modifying the keys. If two
	// transactions were to acquire the Shared lock initially and then attempt
	// to update the keys concurrently, both transactions would get stuck
	// waiting for the other to release its Shared lock and a deadlock would
	// occur. To resolve the deadlock, one of the two transactions would need to
	// be aborted.
	//
	// To avoid this potential deadlock problem, an Upgrade lock can be used in
	// place of a Shared lock. Upgrade locks are not compatible with any other
	// form of locking. As with Shared locks, the lock holder of a Shared lock
	// on a key is only allowed to read from the key while the lock is held.
	// This resolves the deadlock scenario presented above because only one of
	// the transactions would have been able to acquire an Upgrade lock at a
	// time while reading the initial state of the KVs. This means that the
	// Shared-to-Exclusive lock upgrade would never need to wait on another
	// transaction to release its locks.
	//
	// Under pure pessimistic concurrency control, an Upgrade lock is equivalent
	// to an Exclusive lock. However, unlike with Exclusive locks, reads under
	// optimistic concurrency control do not conflict with Upgrade locks. This
	// is because a transaction can only hold an Upgrade lock on keys that it
	// has not yet modified. This improves concurrency between read and write
	// transactions compared to if the writing transaction had immediately
	// acquired an Exclusive lock.
	//
	// The trade-off here is twofold. First, if the Upgrade lock holder does
	// convert its lock on a key to an Exclusive lock after an optimistic read
	// has observed the state of the key, the transaction that performed the
	// optimistic read may be unable to perform a successful read refresh if it
	// attempts to refresh to a timestamp at or past the timestamp of the lock
	// conversion. Second, the optimistic reads permitted while the Upgrade lock
	// is held will bump the timestamp cache. This may result in the Upgrade
	// lock holder being forced to increase its write timestamp when converting
	// to an Exclusive lock, which in turn may force it to restart if its read
	// refresh fails.
	Upgrade Strength = 2
	// Exclusive (X) locks are used by read-write and read-only operations and
	// provide a transaction with exclusive access to a key. When an Exclusive
	// lock is held by a transaction on a given key, no other transaction can
	// read from or write to that key. The lock holder is free to read from and
	// write to the key as frequently as it would like.
	Exclusive Strength = 3
)

var Strength_name = map[int32]string{
	0: "None",
	1: "Shared",
	2: "Upgrade",
	3: "Exclusive",
}
var Strength_value = map[string]int32{
	"None":      0,
	"Shared":    1,
	"Upgrade":   2,
	"Exclusive": 3,
}

func (x Strength) String() string {
	return proto.EnumName(Strength_name, int32(x))
}
func (Strength) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_locking_eecdbf5930525074, []int{0}
}

// Durability represents the different durability properties of a lock acquired
// by a transaction. Durability levels provide varying degrees of survivability,
// often in exchange for the cost of lock acquisition.
type Durability int32

const (
	// Replicated locks are held on at least a quorum of Replicas in a Range.
	// They are slower to acquire and release than Unreplicated locks because
	// updating them requires both cross-node coordination and interaction with
	// durable storage. In exchange, Replicated locks provide a guarantee of
	// survivability across lease transfers, leaseholder crashes, and other
	// forms of failure events. They will remain available as long as their
	// Range remains available and they will never be lost.
	Replicated Durability = 0
	// Unreplicated locks are held only on a single Replica in a Range, which is
	// typically the leaseholder. Unreplicated locks are very fast to acquire
	// and release because they are held in memory or on fast local storage and
	// require no cross-node coordination to update. In exchange, Unreplicated
	// locks provide no guarantee of survivability across lease transfers or
	// leaseholder crashes. They should therefore be thought of as best-effort
	// and should not be relied upon for correctness.
	Unreplicated Durability = 1
)

var Durability_name = map[int32]string{
	0: "Replicated",
	1: "Unreplicated",
}
var Durability_value = map[string]int32{
	"Replicated":   0,
	"Unreplicated": 1,
}

func (x Durability) String() string {
	return proto.EnumName(Durability_name, int32(x))
}
func (Durability) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_locking_eecdbf5930525074, []int{1}
}

// WaitPolicy specifies the behavior of a request when it encounters conflicting
// locks held by other active transactions. The default behavior is to block
// until the conflicting lock is released, but other policies can make sense in
// special situations.
type WaitPolicy int32

const (
	// Block indicates that if a request encounters a conflicting locks held by
	// another active transaction, it should wait for the conflicting lock to be
	// released before proceeding.
	WaitPolicy_Block WaitPolicy = 0
	// Error indicates that if a request encounters a conflicting locks held by
	// another active transaction, it should raise an error instead of blocking.
	// If the request encounters a conflicting lock that was abandoned by an
	// inactive transaction, which is likely due to a transaction coordinator
	// crash, the lock is removed and no error is raised.
	WaitPolicy_Error WaitPolicy = 1
)

var WaitPolicy_name = map[int32]string{
	0: "Block",
	1: "Error",
}
var WaitPolicy_value = map[string]int32{
	"Block": 0,
	"Error": 1,
}

func (x WaitPolicy) String() string {
	return proto.EnumName(WaitPolicy_name, int32(x))
}
func (WaitPolicy) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_locking_eecdbf5930525074, []int{2}
}

func init() {
	proto.RegisterEnum("cockroach.kv.kvserver.concurrency.lock.Strength", Strength_name, Strength_value)
	proto.RegisterEnum("cockroach.kv.kvserver.concurrency.lock.Durability", Durability_name, Durability_value)
	proto.RegisterEnum("cockroach.kv.kvserver.concurrency.lock.WaitPolicy", WaitPolicy_name, WaitPolicy_value)
}

func init() {
	proto.RegisterFile("kv/kvserver/concurrency/lock/locking.proto", fileDescriptor_locking_eecdbf5930525074)
}

var fileDescriptor_locking_eecdbf5930525074 = []byte{
	// 275 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0xcf, 0xbd, 0x4e, 0xc3, 0x30,
	0x10, 0xc0, 0x71, 0x1b, 0x42, 0x69, 0x8f, 0x0f, 0x59, 0x16, 0x53, 0x07, 0x0f, 0x0c, 0x1d, 0x32,
	0x24, 0x03, 0x3c, 0x41, 0x44, 0x57, 0x84, 0xa8, 0x2a, 0x24, 0x36, 0xd7, 0xb1, 0x12, 0x2b, 0x91,
	0x1d, 0x5d, 0x9d, 0x88, 0xbc, 0x01, 0x23, 0xef, 0xc0, 0xcb, 0x74, 0xec, 0xd8, 0x11, 0x92, 0x17,
	0x41, 0x09, 0x20, 0x58, 0x4e, 0x7f, 0x9d, 0xf4, 0x93, 0xee, 0x20, 0x2c, 0x9a, 0xb8, 0x68, 0xb6,
	0x1a, 0x1b, 0x8d, 0xb1, 0x72, 0x56, 0xd5, 0x88, 0xda, 0xaa, 0x36, 0x2e, 0x9d, 0x2a, 0xc6, 0x61,
	0x6c, 0x16, 0x55, 0xe8, 0xbc, 0xe3, 0x0b, 0xe5, 0x54, 0x81, 0x4e, 0xaa, 0x3c, 0x2a, 0x9a, 0xe8,
	0x57, 0x45, 0xff, 0x54, 0x34, 0x80, 0xf9, 0x55, 0xe6, 0x32, 0x37, 0x92, 0x78, 0xa8, 0x6f, 0x1d,
	0x26, 0x30, 0x5d, 0x79, 0xd4, 0x36, 0xf3, 0x39, 0x9f, 0x42, 0x70, 0xef, 0xac, 0x66, 0x84, 0x03,
	0x4c, 0x56, 0xb9, 0x44, 0x9d, 0x32, 0xca, 0xcf, 0xe0, 0x74, 0x5d, 0x65, 0x28, 0x53, 0xcd, 0x8e,
	0xf8, 0x05, 0xcc, 0x96, 0x2f, 0xaa, 0xac, 0xb7, 0xa6, 0xd1, 0xec, 0x78, 0x1e, 0xbc, 0xbe, 0x0b,
	0x12, 0xde, 0x02, 0xdc, 0xd5, 0x28, 0x37, 0xa6, 0x34, 0xbe, 0xe5, 0x97, 0x00, 0x8f, 0xba, 0x2a,
	0x8d, 0x92, 0x5e, 0xa7, 0x8c, 0x70, 0x06, 0xe7, 0x6b, 0x8b, 0x7f, 0x1b, 0xfa, 0xa3, 0xae, 0x01,
	0x9e, 0xa4, 0xf1, 0x0f, 0xae, 0x34, 0xaa, 0xe5, 0x33, 0x38, 0x49, 0x86, 0x33, 0x19, 0x19, 0x72,
	0x89, 0xe8, 0x90, 0xd1, 0x64, 0xb1, 0xfb, 0x14, 0x64, 0xd7, 0x09, 0xba, 0xef, 0x04, 0x3d, 0x74,
	0x82, 0x7e, 0x74, 0x82, 0xbe, 0xf5, 0x82, 0xec, 0x7b, 0x41, 0x0e, 0xbd, 0x20, 0xcf, 0xc1, 0x80,
	0x36, 0x93, 0xf1, 0x99, 0x9b, 0xaf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x5f, 0x18, 0x05, 0x43, 0x38,
	0x01, 0x00, 0x00,
}
