import { useWizardStore } from "./store";
import { stepList } from "./store";

export default function Steps() {
    const { currentStep } = useWizardStore();
    const stepKeys = Object.keys(stepList) as (keyof typeof stepList)[];

    return (
        <ol className="relative text-gray-500 border-s border-gray-200">
            {stepKeys.map((stepKey, index) => {
                const step = stepList[stepKey];
                const isActive = stepKey === currentStep;
                const isPast = stepKeys.indexOf(currentStep) > index;

                return (
                    <li key={stepKey} className="mb-10 ms-6 last:mb-0">
                        <span className={`absolute flex items-center justify-center w-8 h-8 rounded-full -start-4 ring-2 ring-white 
                            ${isPast ? 'bg-green-200' :
                                isActive ? 'bg-gray-100 ring-blue-500' :
                                    'bg-gray-100'}`}>
                            {isPast ? (
                                <svg className="w-3.5 h-3.5 text-green-500" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 16 12">
                                    <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M1 5.917 5.724 10.5 15 1.5" />
                                </svg>
                            ) : (
                                <div className={`${isActive ? 'text-blue-500' : ''}`}>
                                    {step.icon}
                                </div>
                            )}
                        </span>
                        <h3 className={`font-medium leading-tight ${isActive ? 'text-bold text-black' : ''}`}>
                            {step.title}
                        </h3>
                        <p className="text-sm">{step.description}</p>
                    </li>
                );
            })}
        </ol>
    );
}
