SELECT c.id,c.content,cu.name,sc.id,sc.content,scu.name
FROM
    story_t s
        LEFT JOIN comment_t c ON s.id = c.story_id
        LEFT JOIN user_t cu ON c.user_mail = cu.mail
        LEFT JOIN subcomment_t sc ON c.id = sc.comment_id
        LEFT JOIN user_t scu ON sc.user_mail = scu.mail
WHERE
        s.id = 'caf243eb-6358-48a4-be04-3f4765a6d2fa'