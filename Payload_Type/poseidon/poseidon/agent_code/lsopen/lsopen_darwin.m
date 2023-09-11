#import <Foundation/Foundation.h>
#include "lsopen_darwin.h"

// https://github.com/sfsam/Itsycal/blob/11e6e9d265379a610ef103850995e280873f9505/Itsycal/MoLoginItem.m
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"


bool openUsingLSWithh(NSString *path, NSDictionary *env, bool hide, NSArray<NSString*> *arghhhs) {
    FSRef appFSURL;
    OSStatus stat = FSPathMakeRef((const UInt8 *)[path UTF8String], &appFSURL, NULL);
    
    if (stat != errSecSuccess) {
        return false;
    }
    
    LSApplicationParameters appParam;
    appParam.version = 0;
    
    if (hide) {
        appParam.flags = kLSLaunchAndHide;
    } else {
        appParam.flags = kLSLaunchDefaults;
    }

    appParam.application = &appFSURL;
    appParam.argv = (__bridge CFArrayRef) arghhhs;
    //appParam.argv = NULL;
    appParam.environment = (__bridge CFDictionaryRef)env;
    appParam.asyncLaunchRefCon = NULL;
    appParam.initialEvent = NULL;
    CFArrayRef array = (__bridge CFArrayRef)@[];
    stat = LSOpenURLsWithRole(array, kLSRolesAll, NULL, &appParam, NULL, 0);
    if (stat != errSecSuccess) {
        return false;
    }
    return true;
}

int dyldd_inject(char *app, int hide, char * argv[], int argc) {
    @try {
        NSString *appPath = [NSString stringWithCString:app encoding:NSUTF8StringEncoding];
        
        NSDictionary *env = nil;
        
        bool shouldHide = false;
        if (hide == 1) {
            shouldHide = true;
        }

        NSMutableArray *argarray = [NSMutableArray array];
        for (int i = 0; i < argc; i++) {
            NSString *str = [[NSString alloc] initWithCString:argv[i] encoding:NSUTF8StringEncoding];
            [argarray addObject:str];
        }

        NSRange rng = NSMakeRange(2, argc -2);
        NSArray* applicationargs = [argarray subarrayWithRange:rng];
        
        bool success = openUsingLSWithh(appPath, env, shouldHide,applicationargs);
        if (success != true) {
            return -1;
        }
        return 0;
    } @catch (NSException *exception) {
        return -1;
    }
}