export type RemoveProperty = () => string;
type CustomCssVariables = '--scroll-margin-top';
export const setCustomVariable = (property: CustomCssVariables, valueInPx: number) => {
    document.documentElement.style.setProperty(property, `${valueInPx}px`);
    return () => document.documentElement.style.removeProperty(property);
};
