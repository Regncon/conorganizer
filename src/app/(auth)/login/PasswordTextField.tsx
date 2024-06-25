'use client';

import { Visibility, VisibilityOff } from '@mui/icons-material';
import { TextField } from '@mui/material';
import IconButton from '@mui/material/IconButton';
import InputAdornment from '@mui/material/InputAdornment';
import OutlinedInput from '@mui/material/OutlinedInput';
import { useState, MouseEvent, useRef, type RefObject, useEffect, forwardRef } from 'react';

type Props = {
    autoComplete?: string;
    label?: string;
    name?: string;
    error?: boolean;
    helperText?: string;
};

const PasswordTextField = (
    { autoComplete = 'current-password', label = 'passord', name = 'password', error = false, helperText }: Props,
    ref?: RefObject<HTMLInputElement>
) => {
    const [showPassword, setShowPassword] = useState(false);
    const [minCharacterLabel, setMinCharacterLabel] = useState<string>(label);

    const handleClickShowPassword = () => setShowPassword((show) => !show);
    const handleMouseDownPassword = (e: MouseEvent<HTMLButtonElement>) => {
        e.preventDefault();
    };
    const passwordRegExp = /.{6,}/;
    return (
        <TextField
            type={showPassword ? 'text' : 'password'}
            name={name}
            error={error}
            helperText={helperText}
            autoComplete={autoComplete}
            label={minCharacterLabel}
            onFocus={() => {
                setMinCharacterLabel(`${label} minst 6 karakterer`);
            }}
            onBlur={() => {
                setMinCharacterLabel(label);
            }}
            fullWidth
            // required
            InputProps={{
                endAdornment: (
                    <InputAdornment position="end">
                        <IconButton
                            aria-label="toggle password visibility"
                            onClick={handleClickShowPassword}
                            onMouseDown={handleMouseDownPassword}
                            edge="end"
                        >
                            {showPassword ?
                                <VisibilityOff />
                            :   <Visibility />}
                        </IconButton>
                    </InputAdornment>
                ),
            }}
            // inputProps={{
            //     pattern: passwordRegExp.source,
            //     title: 'Minimum antall tegn er 6',
            // }}
        />
    );
};

export default PasswordTextField;
