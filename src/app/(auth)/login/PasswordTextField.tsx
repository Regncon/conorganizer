'use client';

import { Visibility, VisibilityOff } from '@mui/icons-material';
import IconButton from '@mui/material/IconButton';
import InputAdornment from '@mui/material/InputAdornment';
import OutlinedInput from '@mui/material/OutlinedInput';
import { useState, MouseEvent } from 'react';

const PasswordTextField = () => {
	const [showPassword, setShowPassword] = useState(false);

	const handleClickShowPassword = () => setShowPassword((show) => !show);
	const handleMouseDownPassword = (e: MouseEvent<HTMLButtonElement>) => {
		e.preventDefault();
	};

	return (
		<OutlinedInput
			type={showPassword ? 'text' : 'password'}
			name="password"
			autoComplete="current-password"
			label="password"
			required
			endAdornment={
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
			}
		/>
	);
};

export default PasswordTextField;
