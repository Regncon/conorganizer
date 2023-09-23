import { FormControlLabel, FormControlLabelProps, styled, useRadioGroup } from "@mui/material";

interface StyledFormControlLabelProps extends FormControlLabelProps {
    checked: boolean;
}

const StyledFormControlLabel = styled((props: StyledFormControlLabelProps) => <FormControlLabel {...props} />)(
    ({ theme, checked }) => ({
        '.MuiFormControlLabel-label': checked && {
            color: theme.palette.primary.main,
        },
    })
);

const EnrollmentSelector = (props: FormControlLabelProps)  =>{
    const radioGroup = useRadioGroup();
    let checked = false;

    if (radioGroup) {
        checked = radioGroup.value === props.value;
    }

    return <StyledFormControlLabel checked={checked} {...props} />;
}

export default EnrollmentSelector;