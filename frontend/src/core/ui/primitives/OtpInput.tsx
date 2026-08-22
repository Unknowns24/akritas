import React, { useRef, useState, KeyboardEvent, ClipboardEvent } from "react";
import styles from "./OtpInput.module.css";

interface OtpInputProps {
  length?: number;
  value: string;
  onChange: (value: string) => void;
  disabled?: boolean;
}

export const OtpInput: React.FC<OtpInputProps> = ({
  length = 6,
  value,
  onChange,
  disabled = false,
}) => {
  const [activeInput, setActiveInput] = useState<number>(0);
  const inputRefs = useRef<Array<HTMLInputElement | null>>([]);

  const focusInput = (index: number) => {
    const safeIndex = Math.max(0, Math.min(length - 1, index));
    setActiveInput(safeIndex);
    inputRefs.current[safeIndex]?.focus();
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>, index: number) => {
    if (e.key === "Backspace") {
      e.preventDefault();
      const newValue = value.split("");
      if (newValue[index]) {
        newValue[index] = "";
      } else if (index > 0) {
        newValue[index - 1] = "";
        focusInput(index - 1);
      }
      onChange(newValue.join(""));
    } else if (e.key === "ArrowLeft") {
      e.preventDefault();
      focusInput(index - 1);
    } else if (e.key === "ArrowRight") {
      e.preventDefault();
      focusInput(index + 1);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>, index: number) => {
    const val = e.target.value.replace(/[^0-9]/g, ""); // Allow only numbers
    if (!val) return;
    
    // Get only the last char if they typed quickly
    const char = val[val.length - 1];
    const newValue = value.padEnd(length, " ").split("");
    newValue[index] = char;
    onChange(newValue.join("").trim());
    
    if (index < length - 1) {
      focusInput(index + 1);
    }
  };

  const handlePaste = (e: ClipboardEvent<HTMLInputElement>) => {
    e.preventDefault();
    const pastedData = e.clipboardData.getData("text/plain").replace(/[^0-9]/g, "").slice(0, length);
    if (pastedData) {
      onChange(pastedData);
      focusInput(Math.min(pastedData.length, length - 1));
    }
  };

  // Convert string to array of characters, padded with empty strings
  const otpValues = Array.from({ length }, (_, i) => value[i] || "");

  return (
    <div className={styles.container}>
      {otpValues.map((digit, index) => (
        <input
          key={index}
          type="text"
          inputMode="numeric"
          maxLength={2}
          value={digit}
          onChange={(e) => handleChange(e, index)}
          onKeyDown={(e) => handleKeyDown(e, index)}
          onFocus={() => setActiveInput(index)}
          onPaste={handlePaste}
          ref={(el) => { inputRefs.current[index] = el; }}
          disabled={disabled}
          className={`${styles.input} ${activeInput === index ? styles.active : ""}`}
        />
      ))}
    </div>
  );
};
