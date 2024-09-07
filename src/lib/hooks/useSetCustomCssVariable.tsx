import { useEffect, useRef } from 'react';

type BoundingRectProperties = keyof Omit<DOMRect, 'x' | 'y' | 'toJSON'>;
type RemoveProperty = () => void;
type PartialRecord<K extends keyof any, T> = Partial<Record<K, T>>;
type CustomCssVariables = '--scroll-margin-top' | '--app-bar-height';
type PropertyType = PartialRecord<CustomCssVariables, BoundingRectProperties>;

/**
 * Custom hook to dynamically set a CSS variable based on the size of a referenced element.
 * @param {PropertyType} property - The CSS custom property to set.
 * @returns A ref to attach to the DOM element whose size will determine the custom variable.
 */
export const useSetCustomCssVariable = (propertiesMap: PropertyType) => {
    const ref = useRef<HTMLElement | null>(null);

    useEffect(() => {
        let removeProperties: RemoveProperty[] = [];
        const updateCssVariables = () => {
            if (ref.current) {
                const boundingRect = ref.current.getBoundingClientRect();
                removeProperties = Object.entries(propertiesMap).map(([cssVar, boundingProperty]) => {
                    const value = boundingRect[boundingProperty as keyof DOMRect] as number;
                    return setCustomVariable(cssVar, value);
                });
            }
        };
        updateCssVariables();

        const resizeObserver = new ResizeObserver(() => {
            removeProperties.forEach((remove) => remove());
            updateCssVariables();
        });

        if (ref.current) {
            resizeObserver.observe(ref.current);
        }

        return () => {
            resizeObserver.disconnect();
            removeProperties.forEach((remove) => remove());
        };
    }, [propertiesMap]);

    return ref;
};

/**
 * Function to set a custom CSS variable on the document's root element.
 * @param property - The CSS custom property to set.
 * @param valueInPx - The value to set in pixels.
 * @returns A function to remove the property.
 */
const setCustomVariable = (property: string, valueInPx: number): RemoveProperty => {
    document.documentElement.style.setProperty(property, `${valueInPx}px`);
    return () => document.documentElement.style.removeProperty(property);
};
