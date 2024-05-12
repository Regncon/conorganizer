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
};

const PasswordTextField = forwardRef<HTMLInputElement, Props>(
	({ autoComplete = 'current-password', label = 'passord', name = 'password' }, ref) => {
		const [showPassword, setShowPassword] = useState(false);
		const [minCharacterLabel, setMinCharacterLabel] = useState<string>(label);

		const handleClickShowPassword = () => setShowPassword((show) => !show);
		const handleMouseDownPassword = (e: MouseEvent<HTMLButtonElement>) => {
			e.preventDefault();
		};
		const passwordRegExp = /.{8,}/;
		return (
			<TextField
				inputRef={ref}
				type={showPassword ? 'text' : 'password'}
				name={name}
				autoComplete={autoComplete}
				label={minCharacterLabel}
				onFocus={() => {
					setMinCharacterLabel(`${label} minst 8 karakterer`);
				}}
				onBlur={() => {
					setMinCharacterLabel(label);
				}}
				required
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
								:	<Visibility />}
							</IconButton>
						</InputAdornment>
					),
				}}
				inputProps={{
					pattern: passwordRegExp.source,
					title: 'Minimum antall tegn er 8',
				}}
			/>
		);
	}
);

export default PasswordTextField;
