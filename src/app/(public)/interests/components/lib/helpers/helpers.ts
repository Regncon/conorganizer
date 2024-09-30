import type { PoolEvents } from '$app/(public)/components/lib/serverAction';
import { PoolName, type InterestLevel } from '$lib/enums';
import type { Interest, PoolEvent } from '$lib/types';

type ParticipantName = string;
type ParticipantInterestLevel = InterestLevel;
type PoolEventWithInterestLevel = PoolEvent & { interestLevel: ParticipantInterestLevel };
type PoolEventsMapWithInterestLevel = Map<PoolName, PoolEventWithInterestLevel[]>;
type ParticipantPoolEventsMap = Map<ParticipantName, PoolEventsMapWithInterestLevel>;

/**
 * Builds a nested map that associates each participant with their respective pool events.
 *
 * The resulting structure allows easy retrieval of all events a participant is interested in within each pool.
 *
 * **Map Structure:**
 * ```
 * {
 *   "Participant Full Name": {
 *     "Pool Name": [Event, Event, ...],
 *     ...
 *   },
 *   ...
 * }
 * ```
 *
 * @param {Interest[]} interests - An array of participant interests, where each interest links a participant to a specific pool event.
 * @param {PoolEventsMap} poolEventsMap - A map where each key is a pool name and the value is an array of events associated with that pool.
 * @returns {ParticipantPoolEventsMap} A nested map linking each participant to their pool-specific events.
 *
 * @example
 * ```typescript
 * const interests: Interest[] = [
 *   {
 *     participantFirstName: "John",
 *     participantLastName: "Doe",
 *     poolName: "fridayEvening",
 *     poolEventId: "event1"
 *   },
 *   // More interests...
 * ];
 *
 * const poolEventsMap: PoolEventsMap = new Map([
 *   ["fridayEvening", [{ id: "event1", // other properties // }]],
 *   // More pools...
 * ]);
 *
 * const participantMap = buildParticipantPoolEventsMap(interests, poolEventsMap);
 * [...participantMap.entries()].map(([participantName, poolEvents]) => {
 *        return [...poolEvents.entries()].map(([poolName, events]) => {
 *                 return events.map(event => event.id)
 *               })
 *        });
 * ```
 */
export function buildParticipantPoolEventsMap(
    interests: Interest[],
    poolEventsMap: PoolEvents
): ParticipantPoolEventsMap {
    const participantToPoolsMap: ParticipantPoolEventsMap = new Map();

    for (const interest of interests) {
        const { participantFirstName, participantLastName, poolName, poolEventId, interestLevel } = interest;
        const participantFullName = `${participantFirstName} ${participantLastName}`;

        if (!poolName) {
            console.warn(`No pool name specified for participant "${participantFullName}".`);
            continue;
        }

        const eventsInPool = poolEventsMap.get(poolName);
        if (!eventsInPool) {
            console.warn(`Pool "${poolName}" not found for participant "${participantFullName}".`);
            continue;
        }

        const event = { ...eventsInPool.find((e) => e.id === poolEventId), interestLevel } as
            | PoolEventWithInterestLevel
            | undefined;

        if (!event) {
            console.warn(
                `Event ID "${poolEventId}" not found in pool "${poolName}" for participant "${participantFullName}".`
            );
            continue;
        }

        let poolsMap = participantToPoolsMap.get(participantFullName);
        if (!poolsMap) {
            poolsMap = new Map<PoolName, PoolEventWithInterestLevel[]>([
                [PoolName.fridayEvening, []],
                [PoolName.saturdayMorning, []],
                [PoolName.fridayEvening, []],
                [PoolName.saturdayEvening, []],
                [PoolName.sundayMorning, []],
            ]);

            const test = [PoolName.fridayEvening, []];
            participantToPoolsMap.set(participantFullName, poolsMap);
        }

        const events = poolsMap.get(poolName) || [];
        poolsMap.set(
            poolName,
            [...events, event].sort((a, b) => {
                return a.interestLevel > b.interestLevel ? -1 : 1;
            })
        );
    }

    return participantToPoolsMap;
}
