import React from 'react';

type ButtonProps = {
    variant?: 'primary' | 'secondary' | 'danger';
    onClick?: () => void;
    children: React.ReactNode;
    disabled?: boolean;
};

export default function Button({ variant = 'primary', onClick, children, disabled = false }: ButtonProps) {
    const baseStyles = "px-4 py-2 rounded focus:outline-none focus:ring";
    const variantStyles = {
        primary: "bg-blue-500 text-white hover:bg-blue-600 focus:ring-blue-300",
        secondary: "bg-gray-500 text-white hover:bg-gray-600 focus:ring-gray-300",
        danger: "bg-red-500 text-white hover:bg-red-600 focus:ring-red-300",
    };

    const disabledStyles = "opacity-50 cursor-not-allowed";

    return (
        <button
            className={`${baseStyles} ${variantStyles[variant]} ${disabled ? disabledStyles : ''}`}
            onClick={disabled ? undefined : onClick}
            disabled={disabled}
        >
            {children}
        </button>
    );
};
