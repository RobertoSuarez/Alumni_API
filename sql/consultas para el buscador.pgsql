SELECT titulo, area 
    FROM oferta_laborals
    WHERE 
        area in ('frontend', 'Backend') AND
        titulo not like '';