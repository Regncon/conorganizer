import { useMemo } from 'react';
import useReadLocalStorage from '$lib/hooks/useReadLocalStorage';
import type { PoolEvent } from '$lib/types';
import type { GameType } from '$lib/enums';

/**
 * Represents the state of a single filter.
 */
type FilterCriteria = {
    isActive: boolean;
};

/**
 * Aggregates all available filters for PoolEvents, including GameType filters.
 */
type Filters = {
    childFriendly: FilterCriteria;
    adultsOnly: FilterCriteria;
    beginnerFriendly: FilterCriteria;
    lessThanThreeHours: FilterCriteria;
    moreThanSixHours: FilterCriteria;
    possiblyEnglish: FilterCriteria;
    cardGame: FilterCriteria;
    boardGame: FilterCriteria;
    rolePlaying: FilterCriteria;
    other: FilterCriteria;
};

/**
 * Defines GameTypeFilterKeys as 'cardGame' | 'boardGame' | 'rolePlaying' | 'other'.
 */
type GameTypeFilterKeys = 'cardGame' | 'boardGame' | 'rolePlaying' | 'other';

/**
 * Defines BooleanFilterKeys as all keys in Filters excluding GameTypeFilterKeys.
 */
type BooleanFilterKeys = Exclude<keyof Filters, GameTypeFilterKeys>;

/**
 * Defines a type-safe mapping between boolean filter keys and their corresponding PoolEvent properties.
 * Excludes GameType filters.
 */
const BOOLEAN_FILTER_TO_EVENT_PROPERTY_MAP: Record<BooleanFilterKeys, keyof PoolEvent> = {
    childFriendly: 'childFriendly',
    adultsOnly: 'adultsOnly',
    beginnerFriendly: 'beginnerFriendly',
    lessThanThreeHours: 'lessThanThreeHours',
    moreThanSixHours: 'moreThanSixHours',
    possiblyEnglish: 'possiblyEnglish',
} as const;

/**
 * Default filter states when no filters are active or data is missing.
 */
export const DEFAULT_FILTERS: Filters = {
    childFriendly: { isActive: false },
    adultsOnly: { isActive: false },
    beginnerFriendly: { isActive: false },
    lessThanThreeHours: { isActive: false },
    moreThanSixHours: { isActive: false },
    possiblyEnglish: { isActive: false },
    cardGame: { isActive: false },
    boardGame: { isActive: false },
    rolePlaying: { isActive: false },
    other: { isActive: false },
};

/**
 * Type guard to verify if an object conforms to the Filters type.
 * @param obj - The object to verify.
 * @returns True if the object matches the Filters type; otherwise, false.
 */
const isValidFilters = (obj: unknown): obj is Filters => {
    if (typeof obj !== 'object' || obj === null) return false;

    const requiredKeys: Array<keyof Filters> = [
        'childFriendly',
        'adultsOnly',
        'beginnerFriendly',
        'lessThanThreeHours',
        'moreThanSixHours',
        'possiblyEnglish',
        'cardGame',
        'boardGame',
        'rolePlaying',
        'other',
    ];

    return requiredKeys.every((key) => {
        if (!(key in obj)) return false;
        const value = (obj as Record<string, any>)[key];
        return value && typeof value.isActive === 'boolean';
    });
};

/**
 * Determines if a PoolEvent matches any of the active boolean filters.
 * @param event - The PoolEvent to evaluate.
 * @param activeFilterKeys - The keys of active boolean filters.
 * @returns True if the event matches at least one active boolean filter; otherwise, false.
 */
const doesEventMatchBooleanFilters = (event: PoolEvent, activeFilterKeys: Array<BooleanFilterKeys>): boolean => {
    return activeFilterKeys.some(
        (filterKey) =>
            filterKey in BOOLEAN_FILTER_TO_EVENT_PROPERTY_MAP && event[BOOLEAN_FILTER_TO_EVENT_PROPERTY_MAP[filterKey]]
    );
};

/**
 * Determines if a PoolEvent matches any of the active GameType filters.
 * @param event - The PoolEvent to evaluate.
 * @param activeGameTypes - An array of active GameType filters.
 * @returns True if the event's gameType matches at least one active GameType filter; otherwise, false.
 */
const doesEventMatchGameTypeFilters = (event: PoolEvent, activeGameTypes: Array<GameType>): boolean => {
    return activeGameTypes.includes(event.gameType);
};

/**
 * Extracts the keys of all active filters, separated into boolean and GameType filters.
 * @param filters - The current set of filters.
 * @returns An object containing arrays of active boolean filter keys and active GameType filters.
 */
const getActiveFilters = (
    filters: Filters
): { activeBooleanFilters: Array<BooleanFilterKeys>; activeGameTypes: Array<GameType> } => {
    const activeBooleanFilters: Array<BooleanFilterKeys> = [];
    const activeGameTypes: Array<GameType> = [];

    Object.entries(filters).forEach(([key, value]) => {
        if (value.isActive) {
            switch (key as keyof Filters) {
                case 'cardGame':
                case 'boardGame':
                case 'rolePlaying':
                case 'other':
                    activeGameTypes.push(key as GameType);
                    break;
                default:
                    activeBooleanFilters.push(key as BooleanFilterKeys);
            }
        }
    });

    return { activeBooleanFilters, activeGameTypes };
};

/**
 * Determines if a PoolEvent is published.
 * @param event - The PoolEvent to evaluate.
 * @returns True if the event is published; otherwise, false.
 */
const isPublished = (event: PoolEvent): boolean => event.published;

/**
 * Applies active filters to a list of PoolEvents using logical OR across both boolean and GameType filters.
 * @param events - The array of PoolEvents to filter.
 * @param filters - The current set of filters.
 * @returns An array of PoolEvents that match any active filter.
 */
const applyActiveFilters = (events: ReadonlyArray<PoolEvent>, filters: Filters): ReadonlyArray<PoolEvent> => {
    const publishedEvents = events.filter(isPublished);
    const { activeBooleanFilters, activeGameTypes } = getActiveFilters(filters);

    if (activeBooleanFilters.length === 0 && activeGameTypes.length === 0) {
        return publishedEvents;
    }

    return publishedEvents.filter((event) => {
        const matchesBoolean =
            activeBooleanFilters.length > 0 ? doesEventMatchBooleanFilters(event, activeBooleanFilters) : false;
        const matchesGameType =
            activeGameTypes.length > 0 ? doesEventMatchGameTypeFilters(event, activeGameTypes) : false;
        return matchesBoolean || matchesGameType;
    });
};

/**
 * Custom hook to retrieve filtered PoolEvents based on active filters stored in local storage.
 * This hook only reads from localStorage and does not handle writing to it.
 * @param events - The array of PoolEvents to filter.
 * @returns An array of PoolEvents that match the active filters.
 */
export const useFilteredPoolEvents = (events: ReadonlyArray<PoolEvent>): ReadonlyArray<PoolEvent> => {
    const storedFilters = useReadLocalStorage<Filters>('filters');
    const filters: Filters = isValidFilters(storedFilters) ? storedFilters : DEFAULT_FILTERS;

    const filteredEvents = useMemo(() => applyActiveFilters(events, filters), [events, filters]);

    return filteredEvents;
};
