import { useEffect, useState } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import rehypeSanitize from 'rehype-sanitize';
import { fetchUser } from '../services/api';

function formatTime(ts) {
  const d = new Date(ts);
  const hh = String(d.getHours()).padStart(2, '0');
  const mm = String(d.getMinutes()).padStart(2, '0');
  return `${hh}:${mm}`;
}

export default function MessageBubble({ message }) {
  const { uid, headshot, content, time, isSelf } = message;
  const [sender, setSender] = useState({ nickname: `User${uid}`, headshot: '1' });

  useEffect(() => {
    let cancelled = false;
    fetchUser(uid)
      .then((u) => {
        if (!cancelled) setSender(u);
      })
      .catch(() => {});
    return () => {
      cancelled = true;
    };
  }, [uid]);

  const avatarSrc = isSelf ? `/headshot/${headshot}` : `/headshot/${sender.headshot}`;

  return (
    <div className={`msg ${isSelf ? 'msg--self' : 'msg--other'}`}>
      <img
        className="msg__avatar"
        src={avatarSrc}
        alt=""
      />
      <div className="msg__body">
        <div className="msg__meta">
          {!isSelf && <span className="msg__nickname">{sender.nickname}</span>}
          <span className="msg__time">{formatTime(time)}</span>
        </div>
        <div className="msg__content">
          <ReactMarkdown
            remarkPlugins={[remarkGfm]}
            rehypePlugins={[rehypeSanitize]}
            components={{
              a: ({ href, children }) => (
                <a href={href} target="_blank" rel="noopener noreferrer">
                  {children}
                </a>
              ),
            }}
          >
            {content}
          </ReactMarkdown>
        </div>
      </div>
    </div>
  );
}
