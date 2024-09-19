package com.anxcye.runner;

import com.anxcye.constants.RedisConstant;
import com.anxcye.service.ArticleService;
import com.anxcye.utils.RedisCache;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.CommandLineRunner;
import org.springframework.stereotype.Component;

import java.util.Map;

@Component
public class ViewCountRunner implements CommandLineRunner {

    @Autowired
    private ArticleService articleService;

    @Autowired
    private RedisCache redisCache;

    @Override
    public void run(String... args) throws Exception {
        Map<String, Integer> viewCountMap = articleService.getViewCount();
        redisCache.setCacheMap(RedisConstant.ARTICLE_VIEW_COUNT, viewCountMap);
    }
}
