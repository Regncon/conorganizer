import useReadLocalStorage from '$lib/hooks/useReadLocalStorage';
import type { PoolEvent } from '$lib/types';
import { useMemo } from 'react';

/**
 * Represents the state of a single filter.
 */
type FilterCriteria = {
    isActive: boolean;
};

/**
 * Aggregates all available filters for PoolEvents.
 */
type Filters = {
    childFriendly: FilterCriteria;
    adultsOnly: FilterCriteria;
    beginnerFriendly: FilterCriteria;
    lessThanThreeHours: FilterCriteria;
    moreThanSixHours: FilterCriteria;
    possiblyEnglish: FilterCriteria;
};

/**
 * Defines a type-safe mapping between filter keys and PoolEvent properties.
 */
const FILTER_TO_EVENT_PROPERTY_MAP = {
    childFriendly: 'childFriendly',
    adultsOnly: 'adultsOnly',
    beginnerFriendly: 'beginnerFriendly',
    lessThanThreeHours: 'lessThanThreeHours',
    moreThanSixHours: 'moreThanSixHours',
    possiblyEnglish: 'possiblyEnglish',
} as const satisfies Record<keyof Filters, keyof PoolEvent>;

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
};

/**
 * Determines if a PoolEvent is published.
 * @param event - The PoolEvent to evaluate.
 * @returns True if the event is published; otherwise, false.
 */
const isPublished = (event: PoolEvent): boolean => event.published;

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
    ];

    return requiredKeys.every((key) => {
        if (!(key in obj)) return false;
        const value = (obj as Record<string, any>)[key];
        return value && typeof value.isActive === 'boolean';
    });
};

/**
 * Extracts the keys of all active filters.
 * @param filters - The current set of filters.
 * @returns An array of active filter keys.
 */
const getActiveFilterKeys = (filters: Filters): Array<keyof Filters> => {
    return (Object.keys(filters) as Array<keyof Filters>).filter((filterKey) => filters[filterKey].isActive);
};

/**
 * Determines if a PoolEvent matches any of the active filters.
 * @param event - The PoolEvent to evaluate.
 * @param activeFilterKeys - The keys of active filters.
 * @returns True if the event matches at least one active filter; otherwise, false.
 */
const doesEventMatchActiveFilters = (event: PoolEvent, activeFilterKeys: Array<keyof Filters>): boolean => {
    return activeFilterKeys.some((filterKey) => event[FILTER_TO_EVENT_PROPERTY_MAP[filterKey]]);
};

/**
 * Applies active filters to a list of PoolEvents using logical OR.
 * @param events - The array of PoolEvents to filter.
 * @param filters - The current set of filters.
 * @returns An array of PoolEvents that match any active filter.
 */
export const applyActiveFilters = (events: ReadonlyArray<PoolEvent>, filters: Filters): ReadonlyArray<PoolEvent> => {
    const publishedEvents = events.filter(isPublished);
    const activeFilterKeys = getActiveFilterKeys(filters);

    if (activeFilterKeys.length === 0) {
        return publishedEvents;
    }

    return publishedEvents.filter((event) => doesEventMatchActiveFilters(event, activeFilterKeys));
};

/**
 * Custom hook to retrieve filtered PoolEvents based on active filters stored in local storage.
 * @param events - The array of PoolEvents to filter.
 * @returns An array of PoolEvents that match the active filters.
 */
export const useFilteredPoolEvents = (events: ReadonlyArray<PoolEvent>): ReadonlyArray<PoolEvent> => {
    const storedFilters = useReadLocalStorage('filters');
    const filters: Filters = isValidFilters(storedFilters) ? storedFilters : DEFAULT_FILTERS;

    const filteredEvents = useMemo(() => applyActiveFilters(events, filters), [events, filters]);

    return filteredEvents;
};
