import React, { useState, useEffect } from 'react';


interface DashboardProps {
  onLogout: () => void;
}


export default function Dashboard({ onLogout }: DashboardProps) {
  const [isActive, setIsActive] = useState(false);
  
  useEffect(() => {
    setIsActive(true);
  }, []);

  return (
    <div className="relative w-screen h-screen bg-background overflow-hidden">
      <div
        className={`
          fixed transition-all duration-700 ease-in-out z-20
          ${isActive
            ? 'top-4 left-4 md:left-4 text-2xl md:text-3xl'
            : 'top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 text-[clamp(2.5rem,10vw,5rem)]'
          }
          font-mono text-accent tracking-[0.15em]
          after:content-[''] after:inline-block after:align-bottom 
          after:w-[3px] after:h-[1em] after:bg-accent after:ml-[5px] 
          after:animate-cursor-blink
        `}
      >
        WICHTNER
      </div>
      <button onClick={onLogout}>Logout</button>
    </div>
  );
} 