export function formatMessageDate(dateString) {
    const messageDate = new Date(dateString);
    const now = new Date();
    const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
    const yesterday = new Date(today);
    yesterday.setDate(today.getDate() - 1);
    
    // If message is from today, show only time
    if (messageDate >= today) {
        return messageDate.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'});
    }
    
    // If message is from yesterday, show "Yesterday" with time
    if (messageDate >= yesterday) {
        return `Yesterday ${messageDate.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})}`;
    }
    
    // If message is from this year, show month and day with time
    if (messageDate.getFullYear() === now.getFullYear()) {
        return messageDate.toLocaleDateString([], {month: 'short', day: 'numeric'}) + 
               ' ' + messageDate.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'});
    }
    
    // For older messages, show full date with time
    return messageDate.toLocaleDateString([], {year: 'numeric', month: 'short', day: 'numeric'}) + 
           ' ' + messageDate.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'});
}
