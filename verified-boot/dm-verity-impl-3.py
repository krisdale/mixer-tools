#!/usr/bin/env python3
import subprocess
import time
import sys
import os
import urllib.request as request
import json
import fnmatch
import re

disk_image = os.path.abspath(sys.argv[1])
cmd = "losetup -f -P --show {0}".format(disk_image)
print("Executing: " + cmd)
try:
    dev = subprocess.check_output(cmd.split(" ")).decode("utf-8").splitlines()
except Exception:
    raise Exception("{0}: {1}".format(cmd, sys.exc_info()))
print(dev[len(dev) - 1])
boot_num = 1
data_num = 3
hash_num = 4

boot_dev = dev[0] + "p" + str(boot_num)
data_dev = dev[0] + "p" + str(data_num)
hash_dev = dev[0] + "p" + str(hash_num)

verity_name = "vroot1"
initramfs_gen = "initramfs.sh"
boot_gen = "boot.sh"
initramfs_fname = "custom_initramfs.cpio.gz"

subprocess.check_output("rm -rf mnt".split(" "))
subprocess.check_output("mkdir mnt".split(" "))


cmd = "veritysetup --verbose --data-block-size=1024 --hash-block-size=1024 format {0} {1}".format(data_dev, hash_dev)
print("Executing: " + cmd)
try:
    res = subprocess.check_output(cmd.split(" ")).decode("utf-8").splitlines()
except Exception:
    raise Exception("{0}: {1}".format(cmd, sys.exc_info()))

print(res[len(res) - 1])

salt_str = res[len(res) - 3]
root_hash_str = res[len(res) - 2]
print(salt_str)
print(root_hash_str)
salt = salt_str.replace('Salt:            	','')
root_hash = root_hash_str.replace('Root hash:      	','')
print(salt)
print(root_hash)

print("Creating initramfs..")
subprocess.check_output("sh {0}".format(initramfs_gen).split(" "))

print("Generating the init file..")
subprocess.check_output("touch initramfs/init".split(" "))
try:
    outfile = open('initramfs/init', 'w')
    outfile.write("mount -t proc none /proc\n")
    outfile.write("mount -t sysfs none /sys\n")
    outfile.write("mount -t sbin none /sbin\n")
    outfile.write("mount -t bin none /bin\n")
    outfile.write("mount -t lib64 none /lib64\n")
    outfile.write("mount -t devtmpfs none /dev\n")

    outfile.write("veritysetup --verbose --data-block-size=1024 --hash-block-size=1024 create " + verity_name + " /dev/sda" + str(data_num) + " /dev/sda" + str(hash_num) + " " + root_hash + "\n")
    outfile.write("mount /dev/mapper/" + verity_name + " /\n")
    outfile.write("umount /dev\n")
    outfile.write("umount /lib64\n")
    outfile.write("umount /bin\n")
    outfile.write("umount /sbin\n")
    outfile.write("umount /sys\n")
    outfile.write("umount /proc\n")

    outfile.write("exec switch_root / /sbin/init")
    outfile.close()
except IOError:
    print("I/O error")
subprocess.check_output("chmod +x initramfs/init".split(" "))

subprocess.check_output("mount {0} mnt".format(boot_dev).split(" "))
for fname in os.listdir('mnt/loader/entries/'):
    if (fnmatch.fnmatch(fname, 'Clear-*')):
        path = "mnt/loader/entries/" + fname
        try:
            outfile = open(path, 'r')
            content = outfile.read()
            print(content)
            outfile.close()
        except IOError:
            print("I/O error")
        content = re.sub(r"root=* quiet", "quiet", content)
        content = re.sub(r"rw", "", content)
        print(content)

        try:
            outfile = open(path, 'w')
            outfile.write(content)
            outfile.write("initrd EFI/" + initramfs_fname)
            outfile.close()
        except IOError:
            print("I/O error")
subprocess.check_output("umount mnt".split(" "))

print("Updating boot files..")
subprocess.check_output("sh {0} {1} {2}".format(boot_gen, boot_dev, initramfs_fname).split(" "))
#subprocess.check_output("mount {0} mnt".format(boot_dev).split(" "))
##subprocess.check_output("cd initramfs".split(" "))
#os.chdir("initramfs")
#print("Inside dir " + os.getcwd())
#subprocess.check_output("find . -print0 | cpio --null -ov --format=newc | gzip -9 > ../mnt/EFI/{0}".format(initramfs_fname).split(" "))
##subprocess.check_output("cd ..".split(" "))
#os.chdir("../")
#print("Inside dir " + os.getcwd())
#for fname in os.listdir('mnt/loader/entries/'):
#    if (fnmatch.fnmatch(fname, 'Clear-*')):
#        path = "mnt/loader/entries/" + fname
#        print("Writing initramfs config to " + path + " in " + boot_dev)
#        try:
#            outfile = open(path, 'a')
#            outfile.write("initrd EFI/" + initramfs_fname)
#            outfile.close()
#        except IOError:
#            print("I/O error")
#
#
#subprocess.check_output("umount mnt".split(" "))
#subprocess.check_output("rm -rf mnt".split(" "))
#
#cmd = "veritysetup --verbose verify {0} {1} {2}".format(data_dev, hash_dev, root_hash)
#
#print("Executing: " + cmd)
#try:
#   res = subprocess.check_output(cmd.split(" ")).decode("utf-8").splitlines()
#except Exception:
#   raise Exception("{0}: {1}".format(cmd, sys.exc_info()))
#print(res[len(res) - 1])

print("Done!")
